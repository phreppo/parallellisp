(load "list-functions.lisp")

(defun sorted (lst)
    (cond
        ((eq lst nil) t)
        ((eq (length lst) 1) t)
        ((eq (length lst) 2) 
            (<= (car lst) (car (cdr lst))))
        (t
            (and 
                (sorted (first-half lst))
                (sorted (second-half lst))
            )
        )))
