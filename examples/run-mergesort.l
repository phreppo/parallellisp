(load "pmergesort.l")

(write "Merge sorting...")
(time (mergesort llist))

(write "Parallel merge sorting...")
(time (pmergesort llist))

(write "Library merge sorting...")
(time (divide-et-impera mergesort merge llist))

t