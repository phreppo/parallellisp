(defun multl (xs ys)
    (cond 
        ((eq nil xs) nil)
        (t (let ((actual (car xs)))
                (cons (mvl actual ys) (multl (cdr xs) ys))))))

(defun mvl (v lst)
    ;; multiplicate value per list
    (cond 
        ((eq lst nil) nil)
        (t (let ((actual (car lst)))
                (cons (cons v actual) (mvl v (cdr lst)))))))
