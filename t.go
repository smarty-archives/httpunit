package httpunit

type T interface {
	Helper()
	Failed() bool
	Log(args ...interface{})
	Errorf(format string, args ...interface{})
}
