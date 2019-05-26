(defun take (lst n) 
    (cond 
        ((eq n 0) nil) 
        (t (cons (car lst) (take (cdr lst) (1- n))))))
        
(defun drop (lst n) 
    (cond ((eq n 0) lst) 
    (t (drop (cdr lst) (1- n)))))

(defun first-half (lst)
    (take lst (/ (length lst) 2)))

(defun second-half (lst)
    (drop lst (/ (length lst) 2)))