(load "lists.lisp")

(defun search (x lst)
    (cond 
        ((eq lst nil) nil)
        ((eq (length lst) 1) 
            (eq (car lst) x))
        (t (or 
                (search x (first-half lst))
                (search x (second-half lst)))
        )))
