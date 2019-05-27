(load "psearch.lisp")

(write "[SEQ] present first element...")
(time (present 5900 llist)) ; on clisp: stack overflow!
(write "[SEQ] present last element...")
(time (present 9118 llist))
(write "")

(write "[PAR] present first element...")
(time (ppresent 5900 llist))
(write "[PAR] present last element...")
(time (ppresent 9118 llist))
(write "")

(write "[PAR] smart present first element...")
(time (smart-ppresent 5900 llist))
(write "[PAR] smart present last element...")
(time (smart-ppresent 9118 llist))
(write "")

(write "[PAR] genial present first element...")
(time (genial-ppresent 5900 llist))
(write "[PAR] genial present last element...")
(time (genial-ppresent 9118 llist))
(write "")

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
(write "[PAR] library present hand-made closure first element...")
(time (divide-et-impera present-closed or llist))

(setq closed-element 9118)
(write "[PAR] library present hand-made closure last element...")
(time (divide-et-impera present-closed or llist))

(write "[PAR][CLOSURE] library present last element...")
(time (divide-et-impera (present 5900) or llist))