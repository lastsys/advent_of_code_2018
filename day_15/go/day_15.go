package main

import "fmt"

func main() {
	filename := "input.txt"
	part1(filename)
	part2(filename)
}

func part1(filename string) {
	grid := loadData(filename)
	i := 0
	grid.Print()
	for {
		i++
		incomplete := Step(grid)
		if incomplete {
			i--
		}
		fmt.Println(i)
		grid.Print()
		win, race, hp := grid.units.WinCondition()
		if win {
			raceName := "Elf"
			if race == Goblin {
				raceName = "Goblin"
			}
			fmt.Printf("Finished @ %v with %v as winners with total HP = %v\n", i, raceName, hp)
			fmt.Printf("Outcome = %v\n", hp*i)
			break
		}
	}
}

func part2(filename string) {
	elfAttackPower := 0
	for {
		grid := loadData(filename)
		elfAttackPower++
		grid.elfAttackPower = elfAttackPower

		startElfCount := 0
		for _, unit := range grid.units {
			if unit.race == Elf {
				startElfCount++
			}
		}
		fmt.Println(elfAttackPower)

		i := 0
		for {
			i++
			incomplete := Step(grid)
			if incomplete {
				i--
			}
			//grid.Print()
			win, race, hp := grid.units.WinCondition()
			if win {
				if race == Elf {
					finalElfCount := 0
					for _, unit := range grid.units {
						if unit != nil && unit.race == Elf {
							finalElfCount++
						}
					}
					if startElfCount == finalElfCount {
						grid.Print()
						fmt.Printf("Finished @ %v with Elves as winners with total HP = %v and attack power = %v\n",
							i, hp, elfAttackPower)
						fmt.Printf("Outcome = %v\n", i*hp)
						return
					}
				}
				break
			}
		}
	}
}
