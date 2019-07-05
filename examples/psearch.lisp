(load "search.lisp")

(defun ppresent (x lst)
    (cond 
        ((eq lst nil) nil)
        ((eq (length lst) 1) 
            (eq (car lst) x))
        (t {or 
                (ppresent x (first-half lst))
                (ppresent x (second-half lst))}
        )))

(defun lib-ppresent (x myList) 
    (divide-et-impera (present x) or myList))

(setq lib-ppresent-setq 
    (divide-et-impera (ppresent 5900) or))
