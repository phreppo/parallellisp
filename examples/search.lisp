(load "lists.lisp")
(load "list-functions.lisp")

(defun present (x lst)
    (cond 
        ((eq lst nil) nil)
        ((eq (length lst) 1) 
            (eq (car lst) x))
        (t (or 
                (present x (first-half lst))
                (present x (second-half lst)))
        )))

(write "[SEQ] present first element...")
(time (present 5900 llist)) ; on clisp: stack overflow!
(write "[SEQ] present last element...")
(time (present 9118 llist))
(write "")
