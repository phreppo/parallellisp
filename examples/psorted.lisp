(load "parallel.lisp")
(load "lists.lisp")
(load "sorted.lisp")

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