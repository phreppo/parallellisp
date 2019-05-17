(defun fib (n) (cond ((eq n 0) 0) ((eq n 1) 1) (t (+ (fib (- n 1)) (fib (- n 2))))) )

(defun bench (n) (+ (fib n) (fib n) (fib n) (fib n) (fib n) (fib n) (fib n) (fib n) ))

(defun toz (n) (cond ((eq n 0) 0) (t (toz (- n 1))) ) )

(defun inc (n) (cond ((eq n 1000) n) (t (inc (+ n 1))) ))