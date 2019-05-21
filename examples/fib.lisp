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

(defun clisp-bench (n) 
    ; applies fun to n 8 times and then sums these results, sequential
    (+ (fib n) (fib n) (fib n) (fib n) (fib n) (fib n) (fib n) (fib n))) 

(defun clisp-bench (n) 
    ; applies fun to n 8 times and then sums these results, sequential
    (+ (fib n) (fib n) (fib n) (fib n) (fib n) (fib n) (fib n) (fib n)))

(defun p-fib (n) 
    ; computes the fibonacci number of n with the parallel evaluation of subterms
    (cond 
        ((eq n 0) 0) 
        ((eq n 1) 1) 
        (t {+ (fib (- n 1)) (fib (- n 2))})
    ))

{defun pp-fib {n} 
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
(write "
")


(write "seq bench, par fib")
(time (bench p-fib 25))
(write "
")

(write "par bench, seq fib")
(time (p-bench fib 25))
(write "
")

(write "par bench, par fib")
(time (p-bench p-fib 25))
(write "
")

(write "par bench, stupid par fib")
{time {p-bench pp-fib 25}}