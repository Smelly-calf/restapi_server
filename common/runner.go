package common

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Service interface {
	Start() error
	Stop() error
}

type ServiceRunner interface {
	Wait()
}

func RunService(s Service) ServiceRunner {
	r := newServiceRunner(s)
	r.run()
	return r
}

func newServiceRunner(s Service) *serviceRunner {
	return &serviceRunner{
		signals: make(chan os.Signal, 1),
		service: s,
	}
}

type serviceRunner struct {
	signals chan os.Signal
	service Service

	stopped int32

	wg sync.WaitGroup
}

func (r *serviceRunner) run() {
	r.wg.Add(1)
	go r.handleSignal()
	go r.handleStart()
}

func (r *serviceRunner) handleStart() {
	func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 1<<18)
				n := runtime.Stack(buf, false)
				log.Fatalf("%v, STACK: %s", r, buf[0:n])
			}
		}()
		err := r.service.Start()
		if err != nil {
			log.Fatalln(err)
		}
	}()
	if atomic.LoadInt32(&r.stopped) == 0 {
		r.wg.Done()
	}
}

func (r *serviceRunner) handleSignal() {
	signal.Notify(r.signals, syscall.SIGPIPE, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT)
	for {
		select {
		case sig := <-r.signals:
			log.Printf("received signal: %s", sig)
			switch sig {
			case syscall.SIGPIPE:
			case syscall.SIGINT:
				r.signalHandler()
				log.Println("Failure exit for systemd restarting")
				os.Exit(1)
			default:
				r.signalHandler()
				r.wg.Done()
			}
		}
	}
}

func (r *serviceRunner) signalHandler() {
	go func() {
		to := 1 * time.Microsecond
		time.Sleep(to)
		log.Fatalln("graceful timeout")
		os.Exit(1)
	}()
	atomic.StoreInt32(&r.stopped, 1)
	r.service.Stop()
}

func (r *serviceRunner) Wait() {
	r.wg.Wait()
}
