(load "search.lisp")
(load "lists.lisp")

; what i would like to write: (divide-et-impera (ppresent 5900) or llist)
; but (ppresent 5900) is not a value...

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

; when you want another element...
(setq closed-element 11)
(divide-et-impera present-closed or llist)
