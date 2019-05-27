(load "parallel.lisp")
(load "lists.lisp")
(load "list-functions.lisp")

(defun sorted (lst)
    (cond
        ((eq lst nil) t)
        ((eq (length lst) 1) t)
        ((eq (length lst) 2) 
            (< (car lst) (car (cdr lst))))
        (t
            (and 
                (sorted (first-half lst))
                (sorted (second-half lst))
            )
        )))

(setq psorted
    (parallelize
        sorted
        (lambda (lst) 
            (or 
                (eq lst nil)
                (eq (length lst) 1)
                (eq (length lst) 2)
                ))
        first-half
        second-half
        and
    ))