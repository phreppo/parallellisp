(load "parallel.lisp")
(load "sum.lisp")

(defun psumlist (lst)
    (cond
        ((eq lst nil) 0)
        ((eq (length lst) 1) (car lst))
        (t {+   (psumlist (first-half  lst))
                (psumlist (second-half lst))
            })))

(defun smart-sumlist (lst)
    (smart-sumlist-ric lst 1))

(defun smart-sumlist-ric (lst partitions)
    (cond
        ((eq lst nil) 0)
        ((eq (length lst) 1) (car lst))
        ((< partitions ncpu)
            (let ((new-partitions (* partitions 2)))
            {+  (smart-sumlist-ric (first-half  lst) new-partitions)
                (smart-sumlist-ric (second-half lst) new-partitions)
            }))
        (t (+   (sumlist (first-half  lst))
                (sumlist (second-half lst))
            ))))

(write "[PAR] psumlist")
(time (psumlist llist))

(write "[PAR] smart psumlist")
(time (smart-sumlist llist))

(write "[PAR] library divide sum")
(time (divide-et-impera sumlist + llist))