#lang racket
(require math/base)
(require racket/set)
(require racket/generator)

(define (load-data filename)
  (map string->number (file->lines filename #:mode 'text)))

(define (part1 values)
  (printf "Part1: ~a\n" (sum values)))

(define (sum-until-cycle values sum visited)
  (let ([updated-sum (+ sum (values))])
    (if (set-member? visited updated-sum)
        updated-sum
        (sum-until-cycle values updated-sum (set-add visited updated-sum))))) 

(define (part2 values)
  (let ([repeated-values (sequence->repeated-generator values)])
    (let ([cycle-sum (sum-until-cycle repeated-values 0 (set))])    
      (printf "Part2: ~a\n" cycle-sum))))

(let ([values (load-data "../input.txt")])
  (part1 values)
  (part2 values))
