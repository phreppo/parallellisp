(load "search.l")

(defun ppresent (x lst)
    (cond 
        ((eq lst nil) nil)
        ((eq (length lst) 1) 
            (eq (car lst) x))
        (t {or 
                (ppresent x (first-half lst))
                (ppresent x (second-half lst))} 
        )))

(defun smart-ppresent (x lst)
    (smart-ppresent-ric x lst (length lst)))

(defun smart-ppresent-ric (x lst initialSize)
    ; tries to balance the load
    (cond 
        ((eq lst nil) nil)
        ((eq (length lst) 1) 
            (eq (car lst) x))
        ((< (length lst) (/ initialSize ncpu))
            (present x lst))
        (t {or 
                (smart-ppresent-ric x (first-half lst)  initialSize)
                (smart-ppresent-ric x (second-half lst) initialSize)
            } 
        )))

(defun genial-ppresent (x lst)
    (genial-ppresent-ric x lst 1))

(defun genial-ppresent-ric (x lst partitions)
    (cond 
        ((< partitions ncpu)
            ; divide
            (let ((new-partitions (* partitions 2)))
                {or 
                    (genial-ppresent-ric x (first-half lst)  new-partitions)
                    (genial-ppresent-ric x (second-half lst) new-partitions)
                }
            ))
        (t (present x lst))
    ))