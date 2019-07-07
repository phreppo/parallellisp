(load "search.lisp")

(defun psearch (x lst)
    (cond 
        ((eq lst nil) nil)
        ((eq (length lst) 1) 
            (eq (car lst) x))
        (t {or 
                (psearch x (first-half lst))
                (psearch x (second-half lst))}
        )))

(defun lib-psearch (x myList) 
    (divide-et-impera (search x) or myList))

(setq lib-psearch-setq 
    (divide-et-impera (psearch 5900) or))
