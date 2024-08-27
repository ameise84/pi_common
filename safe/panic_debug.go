//go:build debug

package safe

// RecoverPanic debug 模式不抓取panic
func RecoverPanic(PanicHook, string) {
}
