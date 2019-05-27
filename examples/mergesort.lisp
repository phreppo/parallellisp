(load "lists.lisp")
(load "list-functions.lisp")
    
(defun merge (firstList secondList)
  (cond ((not firstList) secondList)
        ((not secondList) firstList)
        ((< (car firstList) (car secondList)) 
            (cons (car firstList) (merge (cdr firstList) secondList)))
        (t 
            (cons (car secondList) (merge firstList (cdr secondList))))))
            
(defun mergesort (lst)
  (cond 
    ((eq (length lst) 1) lst)
    (t (merge 
            (mergesort (take lst (/ (length lst) 2)))
            (mergesort (drop lst (/ (length lst) 2)))))))
        