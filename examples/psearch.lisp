(load "search.lisp")

(defun ppresent (x lst)
    (cond 
        ((eq lst nil) nil)
        ((eq (length lst) 1) 
            (eq (car lst) x))
        (t {or 
                (present x (first-half lst))
                (present x (second-half lst))} 
        )))

(defun smart-ppresent (x lst initialSize)
    (cond 
        ((eq lst nil) nil)
        ((eq (length lst) 1) 
            (eq (car lst) x))
        ((<= (length lst) (/ initialSize ncpu))
            (present x lst))
        (t {or 
                (present x (first-half lst))
                (present x (second-half lst))
            } 
        )))

(write "Parallel present")
(time (ppresent 5900 llist))

(write "Smart present")
(time (ppresent 5900 llist (length llist)))