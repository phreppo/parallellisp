# {parallellisp}

Parallellisp is one *pure funcional* Lisp dialect with some primitives to support parallelism. In particular, parallel argument evaluation is easy handled.

## Features

### Parallelism support

To ask for parallel argument evaluation to the runtime system just use `{}` instead of normal brackets. For example, to calculate the Fibonacci number define:

```lisp
(defun p-fib (n)
    (cond 
        ((eq n 0) 0) 
        ((eq n 1) 1) 
        (t {+ (fib (- n 1)) (fib (- n 2))})
    ))

(p-fib 32)
```

This version is much faster than the implementation with normal brackets. More examples are: [pmergesort](https://github.com/parof/parallellisp/blob/master/examples/pmergesort.lisp), [psearch](https://github.com/parof/parallellisp/blob/master/examples/psearch.lisp), [psorted](https://github.com/parof/parallellisp/blob/master/examples/psorted.lisp) and [psum](https://github.com/parof/parallellisp/blob/master/examples/psum.lisp). Every time a new argument is asked to be evaluated one new [green thread](https://en.wikipedia.org/wiki/Green_threads) is spawned. Spawning too many of them could deteriorate performances, and for this reason one special builtin function is supported: `divide-et-impera`. It tries to optimally allocate the work one the cpus, if one provide the sequential algorithm that works on lists and one way to combine partial results. For example, optimal merge sort would become:

```lisp
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

(defun optimal-mergesort (myList)
    (divide-et-impera mergesort merge myList))

(mergesort oneLargeList) ;; this is slow
(optimal-mergesort oneLargeList) ;; this is fast
```

One can measure performances using the function `time`. Because splitting one data in smaller parts and recombine it is a general pattern that does not apply just to the lists, another primitive has been added to the language: `parallelize`. It acts like `divide-et-impera`, but owrks on generic data. For this reason, one has to provide in addition how to split the data in two partitions and when one base case is faced. For example, the `fib` function would become:

```lisp
(setq pfib 
    (parallelize 
        fib                                 ;; sequential algorithm
        (lambda (n) (or (eq n 0) (eq n 1))) ;; base case detector
        (lambda (n) (- n 2))                ;; split-left
        (lambda (n) (- n 1))                ;; split-right
        +                                   ;; combiner
    ))

(pfib 32)
```
In the last example, with 8 cpus usually one could obtain a speedup of 3x.

## The language

## Install
