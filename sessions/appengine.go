// +build appengine

package sessions

import "os"

// Google App Engine prevents access to the syscall package, which means that we're unable to
// flush our sessions in this case.

func (s *sessionTracker) flushSessionsAndRepeatSignal(shutdown chan<- os.Signal, sig os.Signal) {}

func shutdownSignals() chan os.Signal {
	return make(chan os.Signal)
}
