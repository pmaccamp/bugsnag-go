// +build !appengine

package sessions

import (
	"os"
	"os/signal"
	"syscall"
)

func (s *sessionTracker) flushSessionsAndRepeatSignal(shutdown chan<- os.Signal, sig os.Signal) {
	s.sessionsMutex.Lock()
	defer s.sessionsMutex.Unlock()

	signal.Stop(shutdown)
	if len(s.sessions) > 0 {
		err := s.publisher.publish(s.sessions)
		if err != nil {
			s.config.logf("%v", err)
		}
	}
	syscall.Kill(syscall.Getpid(), sig.(syscall.Signal))
}

func shutdownSignals() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	return c
}
