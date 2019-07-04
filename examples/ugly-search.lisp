(load "search.lisp")
(load "lists.lisp")

; hand-made closure
(defun present-closed (lst)
    (cond 
        ((eq lst nil) nil)
        ((eq (length lst) 1) 
            (eq (car lst) closed-element))
        (t (or 
                (present closed-element (first-half lst))
                (present closed-element (second-half lst)))
        )))

(setq closed-element 5900)
(divide-et-impera present-closed or llist)

;; when you want another element...
(setq closed-element 5900)
(divide-et-impera present-closed or llist)
