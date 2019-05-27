(load "list-functions.lisp")
(load "lists.lisp")


(defun sumlist (lst)
    (cond
        ((eq lst nil) 0)
        ((eq (length lst) 1) (car lst))
        (t (+ 
                (sumlist (first-half  lst))
                (sumlist (second-half lst))
            ))))