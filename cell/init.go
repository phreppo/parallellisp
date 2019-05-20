package cell

import "sync/atomic"

// Init initializes the needed variables. Must be called before using any Lisp structure
func Init() {
	initLanguage()
	initGlobalEnv()
	atomic.AddInt32(&ops, 1)
}
