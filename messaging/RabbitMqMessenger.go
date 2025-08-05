package messaging

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/log"
	"github.com/uoul/go-common/serialization"

	amqp "github.com/rabbitmq/amqp091-go"
)

// -----------------------------------------------------------------------------------
// Type
// -----------------------------------------------------------------------------------

type RabbitMqMessenger struct {
	host     string
	port     uint16
	user     string
	password string

	ctx    context.Context
	logger log.ILogger

	maxRetries    uint
	retryInterval time.Duration
	streamBuffer  uint
	serializer    serialization.ISerializer

	subscriptions map[async.Stream[amqp.Delivery]]subsciption

	addSub    chan subsciptionReq
	removeSub chan async.Stream[amqp.Delivery]
	sendMsg   chan internalMsg
}

type RabbitMqExchange struct {
	Type       string
	Exchange   string
	RoutingKey string
}

type subsciptionReq struct {
	exchange RabbitMqExchange
	sub      async.Stream[amqp.Delivery]
}

type subsciption struct {
	caseIdx  int
	exchange RabbitMqExchange
	queue    *amqp.Queue
	consumer <-chan amqp.Delivery
}

type internalMsg struct {
	Exchange RabbitMqExchange
	Body     []byte
	Retries  uint
}

// -----------------------------------------------------------------------------------
// Public
// -----------------------------------------------------------------------------------

// Publish implements IMessenger.
func (r *RabbitMqMessenger) Publish(topic RabbitMqExchange, msg any) error {
	serializedMsg, err := r.serializer.Marshal(msg)
	if err != nil {
		return err
	}
	r.sendMsg <- internalMsg{
		Exchange: topic,
		Body:     serializedMsg,
	}
	return nil
}

// Subscribe implements IMessenger.
func (r *RabbitMqMessenger) Subscribe(topic RabbitMqExchange) async.Stream[amqp.Delivery] {
	sub := async.NewBufferedStream[amqp.Delivery](r.streamBuffer)
	r.addSub <- subsciptionReq{
		exchange: topic,
		sub:      sub,
	}
	return sub
}

// Unsubscribe implements IMessenger.
func (r *RabbitMqMessenger) Unsubscribe(subsciption async.Stream[amqp.Delivery]) {
	r.removeSub <- subsciption
}

// -----------------------------------------------------------------------------------
// Private
// -----------------------------------------------------------------------------------
func (r *RabbitMqMessenger) run() error {
	// Connect to rabbitmq
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", r.user, r.password, r.host, r.port))
	if err != nil {
		return err
	}
	defer conn.Close()
	// Create channel
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	channelClosed := make(chan *amqp.Error, 1)
	ch.NotifyClose(channelClosed)
	// Init already registered subs
	r.initCurrentSubscriptions(ch)
	// Run
	for {
		// Create select-cases
		cases := r.createSelectCases(channelClosed)
		idx, value, _ := reflect.Select(cases)
		switch idx {
		// Case 0: Parent context done
		case 0:
			return nil
		// Case 1: Connection closed
		case 1:
			err, ok := value.Interface().(*amqp.Error)
			if !ok {
				return fmt.Errorf("failed to cast message from close channel")
			}
			return err
		// Case 2: Send message
		case 2:
			msg, ok := value.Interface().(internalMsg)
			if !ok {
				r.logger.Warning("failed to cast message for sending to rabbitmq")
				continue
			}
			err := ch.ExchangeDeclare(
				msg.Exchange.Exchange, // name
				msg.Exchange.Type,     // type
				false,                 // durable
				true,                  // auto-deleted
				false,                 // internal
				false,                 // no-wait
				nil,                   // arguments
			)
			if err != nil {
				return err
			}
			err = ch.Publish(
				msg.Exchange.Exchange,
				msg.Exchange.RoutingKey,
				false,
				false,
				amqp.Publishing{
					Timestamp: time.Now(),
					Body:      msg.Body,
				},
			)
			if err != nil {
				if msg.Retries < r.maxRetries {
					r.logger.Warningf("try to send message again (%v)...", msg)
					r.sendMsg <- internalMsg{
						Exchange: msg.Exchange,
						Body:     msg.Body,
						Retries:  msg.Retries + 1,
					}
				}
				return fmt.Errorf("failed to publish message to rabbitmq (exchange=%s, routingKey=%s) - %v", msg.Exchange.Exchange, msg.Exchange.RoutingKey, err)
			}
		// Case 3: Add subscription
		case 3:
			req, ok := value.Interface().(subsciptionReq)
			if !ok {
				r.logger.Warning("failed to cast message for subsciption request")
				continue
			}
			// Add subscribtion
			r.subscriptions[req.sub] = subsciption{
				exchange: req.exchange,
				queue:    nil,
				consumer: nil,
			}
			// Bind
			r.declareAndBindQueueForSub(ch, req.sub)
		// Case 4: Remove subsciption
		case 4:
			msg, ok := value.Interface().(async.Stream[amqp.Delivery])
			if !ok {
				r.logger.Warning("failed to cast message for sending to rabbitmq")
				continue
			}
			sub := r.subscriptions[msg]
			ch.QueueDelete(sub.queue.Name, false, false, true)
			delete(r.subscriptions, msg)
		// Case 5-n: Message for subsciption received
		default:
			msg, ok := value.Interface().(amqp.Delivery)
			if !ok {
				r.logger.Warning("failed to cast message for sending to rabbitmq")
				continue
			}
			for k, sub := range r.subscriptions {
				if sub.caseIdx == idx {
					k <- async.ActionResult[amqp.Delivery]{
						Result: msg,
						Error:  nil,
					}
				}
			}
		}
	}
}

func (r *RabbitMqMessenger) initCurrentSubscriptions(ch *amqp.Channel) error {
	for k := range r.subscriptions {
		if err := r.declareAndBindQueueForSub(ch, k); err != nil {
			return err
		}
	}
	return nil
}

func (r *RabbitMqMessenger) declareAndBindQueueForSub(ch *amqp.Channel, key async.Stream[amqp.Delivery]) error {
	sub, exists := r.subscriptions[key]
	if !exists {
		return fmt.Errorf("no subscibtion for key registerd")
	}
	// Declare exchange, if not exists
	err := ch.ExchangeDeclare(
		sub.exchange.Exchange, // name
		sub.exchange.Type,     // type
		false,                 // durable
		true,                  // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return err
	}
	// Declare queue for subscibtion
	q, err := ch.QueueDeclare(
		"",
		false, // Durable
		true,  // AutoDelete
		false, // Exclusive
		false, // No-Wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	// Bind queue
	err = ch.QueueBind(
		q.Name,
		sub.exchange.RoutingKey,
		sub.exchange.Exchange,
		false, // No-Wait
		nil,
	)
	if err != nil {
		return err
	}
	// Create consumer
	consumer, err := ch.Consume(
		q.Name,
		"",
		true,  // Auto-Ack
		false, // Exclusive
		false, // NoLocal
		false, // No-Wait
		nil,
	)
	if err != nil {
		return err
	}
	// Update subscription
	r.subscriptions[key] = subsciption{
		exchange: sub.exchange,
		queue:    &q,
		consumer: consumer,
	}
	return nil
}

func (r *RabbitMqMessenger) createSelectCases(connClosed chan *amqp.Error) []reflect.SelectCase {
	// Create collection
	c := []reflect.SelectCase{}
	// Case 0: Parent context done
	c = append(c, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(r.ctx.Done()),
	})
	// Case 1: Channel closed
	c = append(c, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(connClosed),
	})
	// Case 2: Send message
	c = append(c, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(r.sendMsg),
	})
	// Case 3: Add subscription
	c = append(c, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(r.addSub),
	})
	// Case 4: Remove subsciption
	c = append(c, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(r.removeSub),
	})
	// Case 5-n: Message for subsciption received
	for k, v := range r.subscriptions {
		c = append(c, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(v.consumer),
		})
		r.subscriptions[k] = subsciption{
			caseIdx:  len(c) - 1,
			exchange: v.exchange,
			queue:    v.queue,
			consumer: v.consumer,
		}
	}
	return c
}

// -----------------------------------------------------------------------------------
// Options
// -----------------------------------------------------------------------------------

func WithRabbitMqSerializer(serializer serialization.ISerializer) func(*RabbitMqMessenger) {
	return func(rmm *RabbitMqMessenger) {
		rmm.serializer = serializer
	}
}

func WithRabbitMqRetryInterval(interval time.Duration) func(*RabbitMqMessenger) {
	return func(rmm *RabbitMqMessenger) {
		rmm.retryInterval = interval
	}
}

func WithRabbitMqMaxRetries(retries uint) func(*RabbitMqMessenger) {
	return func(rmm *RabbitMqMessenger) {
		rmm.maxRetries = retries
	}
}

func WithRabbitMqStreamBufferSize(size uint) func(*RabbitMqMessenger) {
	return func(rmm *RabbitMqMessenger) {
		rmm.streamBuffer = size
	}
}

// -----------------------------------------------------------------------------------
// Constructor
// -----------------------------------------------------------------------------------

func NewRabbitMqMessenger(ctx context.Context, logger log.ILogger, host string, port uint16, user string, password string, opts ...func(*RabbitMqMessenger)) IMessenger[RabbitMqExchange, amqp.Delivery] {
	// Init new RabbitMqMessenger
	m := &RabbitMqMessenger{
		ctx:      ctx,
		logger:   logger,
		host:     host,
		port:     port,
		user:     user,
		password: password,

		retryInterval: 10 * time.Second,
		maxRetries:    10,
		streamBuffer:  50,
		serializer:    serialization.NewJSONSerializer(),

		subscriptions: map[async.Stream[amqp.Delivery]]subsciption{},
		sendMsg:       make(chan internalMsg, 50),
		addSub:        make(chan subsciptionReq, 50),
		removeSub:     make(chan async.Stream[amqp.Delivery], 50),
	}
	// Apply options
	for _, o := range opts {
		o(m)
	}
	// Run Messenger
	go func() {
		for {
			select {
			case <-m.ctx.Done():
				return
			default:
				err := m.run()
				if err != nil {
					m.logger.Error(err.Error())
					time.Sleep(m.retryInterval)
				}
			}
		}
	}()
	// Return
	return m
}
