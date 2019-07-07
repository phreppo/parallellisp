(load "search.lisp")
(load "lists.lisp")

; what i would like to write: (divide-et-impera (search 5900) or llist)
; but (ppsearch 5900) is not a value...

; hand-made closure
(defun psearch-closed (lst)
    (cond 
        ((eq lst nil) nil)
        ((eq (length lst) 1) 
            (eq (car lst) closed-element))
        (t (or 
                (psearch closed-element (first-half lst))
                (psearch closed-element (second-half lst)))
        )))

(setq closed-element 5900)
(divide-et-impera psearch-closed or llist)

; when you want another element...
(setq closed-element 11)
(divide-et-impera psearch-closed or llist)
