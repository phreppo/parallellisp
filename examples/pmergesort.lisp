(load "mergesort.lisp")
            
(defun pmergesort (lst)
  (cond 
    ((eq (length lst) 1) lst)
    (t {merge 
            (mergesort (take lst (/ (length lst) 2)))
            (mergesort (drop lst (/ (length lst) 2)))}
        )))

(write "Parallel merge sorting...")
(time (dotimes (n 10) (pmergesort llist)))