(write "running tests...")
((lambda (x y) (x y)) 'car '(1 2 3 4))
((lambda (x y) (+ x y)) 1 2)
(time (+ 1 2 3 4 5 5 6 6 74345 35 34 5 45 245 423545 4 523 454 2352 ))

(defun fib (n) (cond ((eq n 0) 0) ((eq n 1) 1) (t (+ (fib (- n 1)) (fib (- n 2))))))
(defun bench (n) {+ (fib n) (fib n) (fib n) (fib n) (fib n) (fib n) (fib n) (fib n)})
(bench 4)

(defun inc (x) (cond ((eq x 900) t) (t (inc (+ x 1)))))
(inc 100)

(defun toz (n) (cond ((eq n 0) 0) (t (toz (- n 1)))))
(toz 100)

(+ 1 2 (- 3 45))

nil

(cons 1 2)
(car '(1 2 3 4))
(cdr '(1 2 3 4))
(atom nil)
(atom '(1 2 3 4))
(eq 3 (car '(1 2 3 4)))
(let ((x 1) (y 2)) (+ x y))
(eq (let ((x '(7 2 3 4)) (y 2)) (+ (car x) y)) 9)