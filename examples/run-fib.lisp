(load "pfib.lisp")

(write "Sequential fib 32")
(time (fib 32))

(write "Naive pfib 32")
(time (p-fib 32))

(write "Optimal pfib 32")
(time (
    parallelize 
        fib
        (lambda (n) (or (eq n 0) (eq n 1)))
        (lambda (n) (- n 1))
        (lambda (n) (- n 2))
        +
        32 
))

(setq pfib 
    (parallelize 
        fib
        (lambda (n) (or (eq n 0) (eq n 1)))
        (lambda (n) (- n 2))
        (lambda (n) (- n 1))
        +
    ))

(write "Closed pfib 32")
(time (pfib 32))