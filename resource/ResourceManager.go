package resource

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/uoul/go-common/log"
)

type ResourceManager struct {
	resources map[IResource]bool
	finished  chan struct{}
	logger    log.ILogger
}

func NewResourceManager(timeout time.Duration, log log.ILogger) IResourceManager {
	rm := &ResourceManager{
		resources: map[IResource]bool{},
		finished:  make(chan struct{}),
		logger:    log,
	}
	go func() {
		defer close(rm.finished)
		s := make(chan os.Signal, 1)

		// add any other syscalls that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		rm.logger.Info("shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			rm.logger.Infof("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})
		defer timeoutFunc.Stop()

		var wg sync.WaitGroup
		// Do the operations asynchronously to save time
		for r := range rm.resources {
			wg.Add(1)
			go func(r IResource) {
				defer wg.Done()
				if err := r.Close(); err != nil {
					rm.logger.Warningf("clean up failed: %v", err)
				}
			}(r)
		}
		wg.Wait()
		rm.logger.Info("Gracefully shutdown")
	}()
	return rm
}

func (rm *ResourceManager) Register(r IResource) {
	rm.resources[r] = true
}

func (rm *ResourceManager) Unregister(r IResource) {
	delete(rm.resources, r)
}

func (rm *ResourceManager) Wait() {
	<-rm.finished
}
