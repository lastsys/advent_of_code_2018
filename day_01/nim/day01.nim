import strutils
import intsets

proc read_data(filename: string): seq[int] =
  var values = newSeq[int]()
  for line in lines filename:
    values.add(strutils.parseInt(line))
  return values

proc part1(values: seq[int]) =
  var sum = 0
  for value in values:
    sum += value
  echo "Part1: ", sum

proc part2(values: seq[int]) =
  var visited: IntSet
  var v = 0
  while true:
    for _, value in values:
      v += value
      if v in visited:
        echo "Part2: ", v
        return
      else:
        visited.incl(v)

let values = read_data("../input.txt")
part1(values)
part2(values)
