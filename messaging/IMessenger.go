package messaging

import "github.com/uoul/go-common/async"

type IMessenger[K, M any] interface {
	Publish(topic K, msg any) error
	Subscribe(topic K) async.Stream[M]
	Unsubscribe(subsciption async.Stream[M])
}
