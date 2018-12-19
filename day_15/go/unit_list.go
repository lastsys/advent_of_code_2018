package main

type UnitList []*Unit

// Sort according to reading order.
func (ul UnitList) Less(i, j int) bool {
	if ul[i].position[1] < ul[j].position[1] {
		return true
	}
	if ul[i].position[1] == ul[j].position[1] && ul[i].position[0] < ul[j].position[0] {
		return true
	}
	return false
}

func (ul UnitList) Swap(i, j int) {
	ul[i], ul[j] = ul[j], ul[i]
}

func (ul UnitList) Len() int {
	return len(ul)
}

func (ul UnitList) Print() {
	for _, unit := range ul {
		unit.Print()
	}
}

func (ul UnitList) WinCondition() (bool, Race, int) {
	elfCount := 0
	goblinCount := 0
	elfHp := 0
	goblinHp := 0
	for _, unit := range ul {
		if unit != nil {
			switch unit.race {
			case Elf:
				elfCount++
				elfHp += unit.hp
			case Goblin:
				goblinCount++
				goblinHp += unit.hp
			}
		}
	}

	if elfCount > 0 && goblinCount == 0 {
		return true, Elf, elfHp
	}

	if elfCount == 0 && goblinCount > 0 {
		return true, Goblin, goblinHp
	}

	return false, Elf, 0
}
