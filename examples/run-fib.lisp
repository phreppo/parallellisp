(load "pfib.lisp")

(write "seq bench, seq fib")
(time (bench fib 25))
(write "")

(write "seq bench, par fib")
(time (bench p-fib 25))
(write "")

(write "par bench, seq fib")
(time (p-bench fib 25))
(write "")

(write "par bench, par fib")
(time (p-bench p-fib 25))
(write "")

(write "par bench, stupid par fib")
(time (p-bench stupid-p-fib 25))

;; ================================= Just fib =================================

(write "fib 32")
(time (fib 32))

(write "pfib 32")
(time (p-fib 32))

;; it is crazy how keeping the biggest computation on the main thread could change performances
;; ... or maybe no

(write "library fib 32")
(time (
    parallelize 
        fib
        (lambda (n) (or (eq n 0) (eq n 1)))
        (lambda (n) (- n 1))
        (lambda (n) (- n 2))
        +
        32 
))

(write "inverted library fib 32")
(time (
    parallelize 
        fib
        (lambda (n) (or (eq n 0) (eq n 1)))
        (lambda (n) (- n 2))
        (lambda (n) (- n 1))
        +
        32 
))

;; questa cosa e completamente folle
(setq pfib 
    (parallelize 
        fib
        (lambda (n) (or (eq n 0) (eq n 1)))
        (lambda (n) (- n 2))
        (lambda (n) (- n 1))
        +
    ))

(write "closed fib 32")
(time (pfib 32))