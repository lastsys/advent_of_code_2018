import scala.io.Source

object Day02 {
  def loadFile(filename: String): Seq[String] =
    Source
      .fromFile(filename)
      .getLines()
      .map(_.trim)
      .toSeq

  def countPairsAndTriplets(str: String): (Int, Int) = {
    val charCount = str.foldLeft(Map.empty[Char, Int]) { (state, value) =>
      val currentCount = state.getOrElse(value, 0)
      state + (value -> (currentCount + 1))
    }
    charCount.values.foldLeft((0, 0)) { (state, value) =>
      value match {
        case 2 => (state._1 + 1, state._2)
        case 3 => (state._1, state._2 + 1)
        case _ => state
      }
    }
  }

  def part1(serials: Seq[String]): Unit = {
    val (pairs, triplets) = serials.foldLeft((0, 0)) { (count, serial) =>
      countPairsAndTriplets(serial) match {
        case (p, t) if p >= 1 && t >= 1 => (count._1 + 1, count._2 + 1)
        case (p, _) if p >= 1           => (count._1 + 1, count._2)
        case (_, t) if t >= 1           => (count._1, count._2 + 1)
        case _ => count
      }
    }
    val checksum = pairs * triplets
    println(s"Part1: $pairs * $triplets = $checksum")
  }

  def numberOfNonEqualCharacters(s1: String, s2: String): (Int, String) = {
    val (count, rest) = s1.zip(s2).foldLeft((0, List.empty[Char])) { case (state, (c1, c2)) =>
        if (c1 != c2) {
          (state._1 + 1, state._2)
        } else {
          (state._1, c1 +: state._2)
        }
    }
    (count, rest.reverse.mkString)
  }

  def part2(serials: Seq[String]): Unit = {
    serials.zipWithIndex.foreach { case (serial1, i) =>
      serials.drop(i).foreach { serial2 =>
        val (count, rest) = numberOfNonEqualCharacters(serial1, serial2)
        if (count == 1) {
          println(s"$count : $rest")
        }
      }
    }
  }

  def main(args: Array[String]): Unit = {
    val serials = loadFile("../input.txt")
    part1(serials)
    part2(serials)
  }
}
