(load "psorted.l")
(load "pmergesort.l")

(write "[SEQ] checking for llist sorted")
(time (sorted llist))

(write "[PAR] checking for llist sorted")
(time (psorted llist))

(setq sortedList (mergesort llist))

(write "[SEQ] checking for sorted list")
(time (sorted sortedList))

(write "[PAR] checking for sorted list")
(time (psorted sortedList))