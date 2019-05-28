package cell

func Eval(c Cell) EvalResult {
	return eval(c, emptyEnv())
}

type EvalResult struct {
	Cell Cell
	Err  error
}

type EvalError struct {
	Err string
}

func (e EvalError) Error() string {
	return e.Err
}

// Init initializes the needed variables. Must be called before using any lisp structure
func Init() {
	initLanguage()
	initglobalEnv()
}
