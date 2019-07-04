(defun fib (n) 
    ; computes the fibonacci number of n
    (cond 
        ((eq n 0) 0) 
        ((eq n 1) 1) 
        (t (+ (fib (- n 1)) (fib (- n 2))))
    ))

(defun bench (n) 
    ; applies fun to n 8 times and then sums these results
    (+ (fib n) (fib n) (fib n) (fib n) (fib n) (fib n) (fib n) (fib n))) 