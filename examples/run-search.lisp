(load "psearch.lisp")

(write "[SEQ] search first element...")
(time (search 5900 llist)) ; on clisp: stack overflow!
(write "[SEQ] search last element...")
(time (search 9118 llist))
(write "")

(write "[PAR] search first element...")
(time (psearch 5900 llist))
(write "[PAR] search last element...")
(time (psearch 9118 llist))
(write "")

; hand-made closure
(defun search-closed (lst)
    (cond 
        ((eq lst nil) nil)
        ((eq (length lst) 1) 
            (eq (car lst) closed-element))
        (t (or 
                (search closed-element (first-half lst))
                (search closed-element (second-half lst)))
        )))

(setq closed-element 5900)
(write "[PAR] library search hand-made closure first element...")
(time (divide-et-impera search-closed or llist))

(setq closed-element 9118)
(write "[PAR] library search hand-made closure last element...")
(time (divide-et-impera search-closed or llist))

(write "[PAR][CLOSURE] library search last element...")
(time (divide-et-impera (search 5900) or llist))