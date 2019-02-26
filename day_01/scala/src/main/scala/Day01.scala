import scala.annotation.tailrec
import scala.io.Source

object Day01 {
  def loadFile(filename: String): Array[Int] =
    Source
      .fromFile(filename)
      .getLines
      .map(_.trim)
      .map(_.toInt)
      .toArray

  def part1(values: Array[Int]): Unit = {
    println(s"Answer part 1: ${values.sum}")
  }

  def part2(values: Array[Int]): Unit = {
    // Loop array infinitely.
    val infiniteValues = Stream.continually(values.toStream).flatten

    @tailrec
    def sumUntilCycle(values: Stream[Int], sum: Int = 0, visited: Set[Int] = Set.empty): Int = {
      val updatedSum = sum + values.head
      if (visited.contains(updatedSum)) {
        return updatedSum
      }
      sumUntilCycle(values.tail, updatedSum, visited + updatedSum)
    }

    println(s"Answer part 2: ${sumUntilCycle(infiniteValues)}")
  }

  def main(args: Array[String]): Unit = {
    val values = loadFile("../input.txt")
    part1(values)
    part2(values)
  }
}
