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

(defun parallelize (sequential-algorithm is-base-case split-left split-right combinator  generic-data)
    (parallelize-ric  1 sequential-algorithm is-base-case split-left split-right combinator  generic-data))

(defun parallelize-ric (partitions sequential-algorithm is-base-case split-left split-right combinator generic-data)
    (cond
        ((is-base-case generic-data) (sequential-algorithm generic-data))
        ((< partitions ncpu)
            (let ((new-partitions (* partitions 2)))
            {combinator 
                (parallelize-ric new-partitions sequential-algorithm is-base-case split-left split-right combinator (split-left generic-data))
                (parallelize-ric new-partitions sequential-algorithm is-base-case split-left split-right combinator (split-right generic-data))
            }))
        (t (combinator 
                (sequential-algorithm (split-left generic-data))
                (sequential-algorithm (split-right generic-data))
            ))))

