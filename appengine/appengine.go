package bugsnagappengine

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bugsnag/bugsnag-go"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	aelog "google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/user"
)

func init() {
	bugsnag.OnBeforeNotify(appengineMiddleware)
}

func appengineMiddleware(event *bugsnag.Event, config *bugsnag.Configuration) (err error) {
	ctx, err := findContext(event)
	if err != nil {
		return err
	}

	// You can only use the builtin HTTP library if you pay for appengine,
	// so we use the appengine urlfetch service instead.
	config.Transport = &urlfetch.Transport{Context: ctx, AllowInvalidServerCertificate: false}
	config.Logger = findLogger(ctx, config.Logger)
	config.ReleaseStage = findReleaseStage(config.ReleaseStage)

	event.User = findUser(ctx, event.User)

	return nil
}

type logger interface {
	Printf(format string, v ...interface{})
}

func findUser(ctx context.Context, bUser *bugsnag.User) *bugsnag.User {
	if u := user.Current(ctx); u != nil && bUser != nil {
		return &bugsnag.User{Id: u.ID, Email: u.Email}
	}
	return bUser
}

func findLogger(ctx context.Context, l logger) logger {
	if configuredLogger, ok := l.(*log.Logger); ok {
		return log.New(appengineWriter{ctx}, configuredLogger.Prefix(), configuredLogger.Flags())
	}
	return log.New(appengineWriter{ctx}, log.Prefix(), log.Flags())
}

func findReleaseStage(rs string) string {
	if rs != "" {
		return rs
	}

	if appengine.IsDevAppServer() {
		return "development"
	}
	return "production"

}

func findContext(event *bugsnag.Event) (context.Context, error) {
	var ctx context.Context

	for _, datum := range event.RawData {
		if r, ok := datum.(*http.Request); ok {
			ctx = appengine.NewContext(r)
			break
		} else if context, ok := datum.(context.Context); ok {
			ctx = context
			break
		}
	}

	var err error
	if ctx == nil {
		err = fmt.Errorf("No appengine context given")
	}
	return ctx, err

}

// Create a custom writer so we can set up an internal logger for Bugsnag
type appengineWriter struct {
	ctx context.Context
}

func (w appengineWriter) Write(b []byte) (int, error) {
	aelog.Errorf(w.ctx, string(b))
	return len(b), nil
}
