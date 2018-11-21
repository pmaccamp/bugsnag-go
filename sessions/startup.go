package sessions

import (
	"context"
	"net/http"
	"os"
)

// This is a copy of panicwrap (panicwrap.DEFAULT_COOKIE_KEY). We cannot
// directly reference this constant because panicwrap imports the 'syscall'
// package, which will prevent app engine applications from being built.
const applicationProcessKey = "cccf35992f8f3cd8d1d28f0109dd953e26664531"

// SendStartupSession is called by Bugsnag on startup, which will send a
// session to Bugsnag and return a context to represent the session of the main
// goroutine. This is the session associated with any fatal panics that are
// caught by panicwrap.
func SendStartupSession(config *SessionTrackingConfiguration) context.Context {
	ctx := context.Background()
	session := newSession()
	if !config.IsAutoCaptureSessions() || isApplicationProcess() {
		return ctx
	}
	publisher := &publisher{
		config: config,
		client: &http.Client{Transport: config.Transport},
	}
	go publisher.publish([]*Session{session})
	return context.WithValue(ctx, contextSessionKey, session)
}

// Checks to see if this is the application process, as opposed to the process
// that monitors for panics
func isApplicationProcess() bool {
	// Application process is run first, and this will only have been set when
	// the monitoring process runs
	return "" == os.Getenv(applicationProcessKey)
}
