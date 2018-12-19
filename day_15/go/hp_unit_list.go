package main

type HpUnitList []*Unit

// Sort according to reading order.
func (hp HpUnitList) Less(i, j int) bool {
	// First check hitpoint.
	if hp[i].hp < hp[j].hp {
		return true
	}
	// If a tie, the reading order wins.
	if hp[i].hp == hp[j].hp {
		if hp[i].position[1] < hp[j].position[1] {
			return true
		}
		if hp[i].position[1] == hp[j].position[1] && hp[i].position[0] < hp[j].position[0] {
			return true
		}
	}
	return false
}

func (hp HpUnitList) Swap(i, j int) {
	hp[i], hp[j] = hp[j], hp[i]
}

func (hp HpUnitList) Len() int {
	return len(hp)
}
