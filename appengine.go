// +build appengine

package bugsnag

// panicwrap is not supported on Google App Engine due to restrictions in the
// available core packages that Google permits access to, in particular,
// 'syscall'.
func defaultPanicHandler() {}
