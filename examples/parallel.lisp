(load "list-functions.lisp")

(defun divide-et-impera (sequential-algorithm combinator lst)
    (divide-et-impera-ric 1 sequential-algorithm combinator lst))

(defun divide-et-impera-ric (partitions sequential-algorithm combinator lst)
    (cond
        ((eq lst nil)        (sequential-algorithm lst))
        ((eq (length lst) 1) (sequential-algorithm lst))
        ((< partitions ncpu)
            (let ((new-partitions (* partitions 2)))
            {combinator 
                (divide-et-impera-ric new-partitions sequential-algorithm combinator (first-half  lst))
                (divide-et-impera-ric new-partitions sequential-algorithm combinator (second-half lst))
            }))
        (t (combinator 
                (sequential-algorithm (first-half  lst))
                (sequential-algorithm (second-half lst))
            ))))

(defun parallelize (sequential-algorithm is-base-case split-left split-right combinator  generic-data)
    ;; contract: split left and split right must be interchangeable and must split
    ;; the data in two parts of equal dimension  
    (parallelize-ric  1 sequential-algorithm is-base-case split-left split-right combinator  generic-data))

(defun parallelize-ric (partitions sequential-algorithm is-base-case split-left split-right combinator generic-data)
    (cond
        ((is-base-case generic-data) (sequential-algorithm generic-data))
        ((< partitions ncpu)
            (let ((new-partitions (* partitions 2))) 
            {combinator 
                (parallelize-ric new-partitions sequential-algorithm is-base-case split-right split-left combinator (split-left generic-data))
                (parallelize-ric new-partitions sequential-algorithm is-base-case split-right split-left combinator (split-right generic-data))
            }))
        (t (combinator 
                (sequential-algorithm (split-left generic-data))
                (sequential-algorithm (split-right generic-data))
            ))))

