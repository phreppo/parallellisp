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

(defun parallelize (generic-data is-base-case split-one split-two combinator sequential-algorithm)
    (parallelize-ric generic-data 1 is-base-case split-one split-two combinator sequential-algorithm))

(defun parallelize-ric (generic-data partitions is-base-case split-one split-two combinator sequential-algorithm)
    (cond
        ((is-base-case generic-data) (sequential-algorithm generic-data))
        ((< partitions ncpu)
            (let ((new-partitions (* partitions 2)))
            {combinator 
                (parallelize-ric (split-one generic-data) new-partitions is-base-case split-one split-two combinator sequential-algorithm)
                (parallelize-ric (split-two generic-data) new-partitions is-base-case split-one split-two combinator sequential-algorithm)
            }))
        (t (combinator 
                (sequential-algorithm (split-one generic-data))
                (sequential-algorithm (split-two generic-data))
            ))))

