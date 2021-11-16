package godzilla

func catchPanic(f func()) (recv interface{}) {
	defer func() {
		recv = recover()
	}()

	f()

	return
}

type testRouter struct {
	path     string
	conflict bool
}
