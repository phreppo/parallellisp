(load "lists.lisp")

(defun take (lst n) 
    (cond 
        ((eq n 0) nil) 
        (t (cons (car lst) (take (cdr lst) (- n 1))))))
        
(defun drop (lst n) 
    (cond ((eq n 0) lst) 
    (t (drop (cdr lst) (- n 1)))))
    
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
            
(write "Merge sorting...")
(time (dotimes (n 10) (mergesort llist)))
