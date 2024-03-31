package ctry

// Try 只捕获任意panic 不做任何处理 test
func Try(f func()) (r interface{}) {
	defer func() {
		r = recover()
	}()
	f()
	return
}
