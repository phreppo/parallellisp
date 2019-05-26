(load "list-functions.lisp")

(defun go (lst combinator sequential-algorithm)
    (go-ric lst 1 combinator sequential-algorithm))

(defun go-ric (lst partitions combinator sequential-algorithm)
    (cond
        ((eq lst nil) (sequential-algorithm lst))
        ((eq (length lst) 1) (sequential-algorithm lst))
        ((< partitions ncpu)
            (let ((new-partitions (* partitions 2)))
            {combinator 
                (go-ric (first-half  lst) new-partitions combinator sequential-algorithm)
                (go-ric (second-half lst) new-partitions combinator sequential-algorithm)
            }))
        (t (combinator 
                (sequential-algorithm (first-half  lst))
                (sequential-algorithm (second-half lst))
            ))))

