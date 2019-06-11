(load "psum.l")

(write "[SEQ] sumlist")
(time (sumlist llist))

(write "[PAR] psumlist")
(time (psumlist llist))

(write "[PAR] smart psumlist")
(time (smart-sumlist llist))

(write "[PAR] library divide sum")
(time (divide-et-impera sumlist + llist))