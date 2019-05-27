(load "parallel.lisp")

(defun fib (n) 
    ; computes the fibonacci number of n in the sequential way
    (cond 
        ((eq n 0) 0) 
        ((eq n 1) 1) 
        (t (+ (fib (- n 1)) (fib (- n 2))))
    ))

(defun bench (fun n) 
    ; applies fun to n 8 times and then sums these results, sequential
    (+ (fun n) (fun n) (fun n) (fun n) (fun n) (fun n) (fun n) (fun n)))

(defun p-fib (n) 
    ; computes the fibonacci number of n with the parallel evaluation of subterms
    (cond 
        ((eq n 0) 0) 
        ((eq n 1) 1) 
        (t {+ (fib (- n 1)) (fib (- n 2))})
    ))

{defun stupid-p-fib {n} 
    ; computes the fibonacci number of n with the parallel evaluation of subterms
    {cond 
        {{eq n 0} 0} 
        {{eq n 1} 1} 
        {t {+ {fib {- n 1}} {fib {- n 2}}}}
    }}

(defun p-bench (fun n) 
    ; applies fun to n 8 times and then sums these results, parallel
    {+ (fun n) (fun n) (fun n) (fun n) (fun n) (fun n) (fun n) (fun n)})
    
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
{time {p-bench stupid-p-fib 25}}

(write "fib 32")
(time (fib 32))

;; it is crazy how keeping the biggest computation on the main thread could change performances
;; ... or maybe no

(write "library fib 32")
(time (
    parallelize 
        32 
        (lambda (n) (or (eq n 0) (eq n 1)))
        (lambda (n) (- n 1))
        (lambda (n) (- n 2))
        +
        fib
))

(write "inverted library fib 32")
(time (
    parallelize 
        32 
        (lambda (n) (or (eq n 0) (eq n 1)))
        (lambda (n) (- n 2))
        (lambda (n) (- n 1))
        +
        fib
))

