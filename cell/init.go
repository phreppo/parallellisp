package cell

// Init initializes the needed variables. Must be called before using any Lisp structure
func Init() {
	initLanguage()
	initGlobalEnv()
}
