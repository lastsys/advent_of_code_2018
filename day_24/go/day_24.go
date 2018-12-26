package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Trait int

const (
	Radiation Trait = iota
	Bludgeoning
	Fire
	Cold
	Slashing
)

type GroupType int

const (
	ImmuneType GroupType = iota
	InfectionType
)

func (g GroupType) String() string {
	if g == ImmuneType {
		return "Immune System"
	}
	return "Infection"
}

func CreateTrait(text string) Trait {
	switch strings.ToLower(text) {
	case "radiation":
		return Radiation
	case "bludgeoning":
		return Bludgeoning
	case "fire":
		return Fire
	case "cold":
		return Cold
	case "slashing":
		return Slashing
	}
	log.Fatalf("Failed to create trait from %v.", text)
	return 0
}

type Group struct {
	Id            int
	Type          GroupType
	UnitHitPoints int
	UnitCount     int
	Weakness      map[Trait]bool
	Immune        map[Trait]bool
	Initiative    int
	AttackDamage  int
	AttackTrait   Trait
}

func (g *Group) IsImmune(trait Trait) bool {
	_, ok := g.Immune[trait]
	return ok
}

func (g *Group) IsWeakness(trait Trait) bool {
	_, ok := g.Weakness[trait]
	return ok
}

func (g *Group) EffectivePower() int {
	return g.UnitCount * g.AttackDamage
}

type GroupList []*Group

type System struct {
	ImmuneSystem GroupList
	Infection    GroupList
}

func NewSystem() *System {
	return &System{
		make(GroupList, 0),
		make(GroupList, 0),
	}
}

func (s *System) Print() {
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println("Immune System:")
	for _, g := range s.ImmuneSystem {
		fmt.Printf("Group %d (%s) contains %d units\n", g.Id, g.Type.String(), g.UnitCount)
	}
	fmt.Println("Infection:")
	for _, g := range s.Infection {
		fmt.Printf("Group %d (%s) contains %d units\n", g.Id, g.Type.String(), g.UnitCount)
	}
}

func (s *System) TargetSelectionOrder() GroupList {
	v := make(GroupList, 0)
	v = append(v, s.Infection...)
	v = append(v, s.ImmuneSystem...)

	selectOrder := func(i, j int) bool {
		if v[i].EffectivePower() > v[j].EffectivePower() {
			return true
		}
		if v[i].EffectivePower() == v[j].EffectivePower() &&
			v[i].Initiative > v[j].Initiative {
			return true
		}
		return false
	}

	sort.Slice(v, selectOrder)

	return v
}

func Damage(attacking *Group, defending *Group) int {
	damage := attacking.EffectivePower()
	if _, ok := defending.Immune[attacking.AttackTrait]; ok {
		damage = 0
	}
	if _, ok := defending.Weakness[attacking.AttackTrait]; ok {
		damage *= 2
	}
	return damage
}

func (s *System) MaxDamage(attackingGroup *Group, defendingGroups map[*Group]bool) *Group {
	targetGroups := s.ImmuneSystem
	if attackingGroup.Type == ImmuneType {
		targetGroups = s.Infection
	}
	targetGroupsCopy := make(GroupList, 0, len(targetGroups))
	for _, g := range targetGroups {
		// Check if already defending.
		if _, ok := defendingGroups[g]; !ok {
			// Check if unit is immune to attack.
			if _, ok := g.Immune[attackingGroup.AttackTrait]; !ok {
				targetGroupsCopy = append(targetGroupsCopy, g)
			}
		}
	}
	f := func(i, j int) bool {
		damage1 := Damage(attackingGroup, targetGroupsCopy[i])
		damage2 := Damage(attackingGroup, targetGroupsCopy[j])
		if damage1 > damage2 {
			return true
		}
		if damage1 == damage2 {
			if targetGroupsCopy[i].EffectivePower() > targetGroupsCopy[j].EffectivePower() {
				return true
			}
			if targetGroupsCopy[i].EffectivePower() == targetGroupsCopy[j].EffectivePower() {
				if targetGroupsCopy[i].Initiative > targetGroupsCopy[j].Initiative {
					return true
				}
			}
		}
		return false
	}
	sort.Slice(targetGroupsCopy, f)
	if len(targetGroupsCopy) == 0 || Damage(attackingGroup, targetGroupsCopy[0]) == 0 {
		return nil
	}
	return targetGroupsCopy[0]
}

func (s *System) RemoveKilledGroups() {
	immune := make(GroupList, 0)
	for _, g := range s.ImmuneSystem {
		if g.UnitCount > 0 {
			immune = append(immune, g)
		}
	}
	infection := make(GroupList, 0)
	for _, g := range s.Infection {
		if g.UnitCount > 0 {
			infection = append(infection, g)
		}
	}
	s.ImmuneSystem = immune
	s.Infection = infection
}

func (s *System) Boost(boost int) {
	for _, g := range s.ImmuneSystem {
		g.AttackDamage += boost
	}
}

func DoAttack(attacking *Group, defending *Group) (int, int) {
	damage := Damage(attacking, defending)
	killedUnits := damage / defending.UnitHitPoints
	if killedUnits > defending.UnitCount {
		return damage, defending.UnitCount
	}
	return damage, killedUnits
}

func loadScenario(filename string) *System {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	pattern := regexp.MustCompile(`(?P<units>\d+) units each with (?P<hitpoints>\d+) hit points (?:\((?P<traits>.+)\))?\s?with an attack that does (?P<damage>\d+) (?P<attacktrait>\w+) damage at initiative (?P<initiative>\d+)`)

	scanner := bufio.NewScanner(f)

	system := NewSystem()

	type Mode int
	const (
		immune Mode = iota
		infection
	)

	currentSystem := immune
	id := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "Immune System:" {
			currentSystem = immune
			id = 0
			continue
		}
		if line == "Infection:" {
			currentSystem = infection
			id = 0
			continue
		}
		if matches := pattern.FindStringSubmatch(line); len(matches) > 0 {
			group := &Group{}
			id++
			group.Id = id
			switch currentSystem {
			case immune:
				group.Type = ImmuneType
			case infection:
				group.Type = InfectionType
			}
			group.Weakness = make(map[Trait]bool)
			group.Immune = make(map[Trait]bool)
			names := pattern.SubexpNames()
			for i := 1; i < len(names); i++ {
				switch names[i] {
				case "units":
					if units, err := strconv.Atoi(matches[i]); err != nil {
						log.Fatal(err)
					} else {
						group.UnitCount = units
					}
				case "hitpoints":
					if hp, err := strconv.Atoi(matches[i]); err != nil {
						log.Fatal(err)
					} else {
						group.UnitHitPoints = hp
					}
				case "traits":
					weakTraits, immuneTraits := parseTraits(matches[i])
					for _, trait := range weakTraits {
						group.Weakness[trait] = true
					}
					for _, trait := range immuneTraits {
						group.Immune[trait] = true
					}
				case "damage":
					if damage, err := strconv.Atoi(matches[i]); err != nil {
						log.Fatal(err)
					} else {
						group.AttackDamage = damage
					}
				case "attacktrait":
					group.AttackTrait = CreateTrait(matches[i])
				case "initiative":
					if initiative, err := strconv.Atoi(matches[i]); err != nil {
						log.Fatal(err)
					} else {
						group.Initiative = initiative
					}
				}
			}
			switch currentSystem {
			case immune:
				system.ImmuneSystem = append(system.ImmuneSystem, group)
			case infection:
				system.Infection = append(system.Infection, group)
			}
		}
	}
	return system
}

func parseTraits(traits string) ([]Trait, []Trait) {
	pattern := regexp.MustCompile(`(?P<type>weak|immune) to (?P<trait>[^;\n]+)`)
	matches := pattern.FindAllStringSubmatch(traits, -1)
	weak := make([]Trait, 0)
	immune := make([]Trait, 0)
	for _, m := range matches {
		t := strings.Split(m[2], ", ")
		switch m[1] {
		case "weak":
			for _, trait := range t {
				weak = append(weak, CreateTrait(strings.TrimSpace(trait)))
			}
		case "immune":
			for _, trait := range t {
				immune = append(immune, CreateTrait(strings.TrimSpace(trait)))
			}
		}
	}
	return weak, immune
}

func FullAttack(system *System, print bool) {
	selectOrder := system.TargetSelectionOrder()
	defendingImmuneGroups := make(map[*Group]bool)
	defendingInfectionGroups := make(map[*Group]bool)
	attack := make([][2]*Group, 0)
	for _, g := range selectOrder {
		defending := defendingImmuneGroups
		if g.Type == ImmuneType {
			defending = defendingInfectionGroups
		}
		if _, ok := defending[g]; !ok {
			target := system.MaxDamage(g, defending)
			defending[target] = true
			attack = append(attack, [2]*Group{g, target})
		}
	}
	if print {
		fmt.Println()
	}
	attackSorter := func(i, j int) bool {
		if attack[i][0].Initiative > attack[j][0].Initiative {
			return true
		}
		return false
	}
	sort.Slice(attack, attackSorter)
	for _, g := range attack {
		if g[1] != nil {
			damage, killedUnits := DoAttack(g[0], g[1])
			if print {
				fmt.Printf("Attack %v (%s) - Defend %v (%s) - killed %d units out of %d with %d damage.\n",
					g[0].Id, g[0].Type.String(), g[1].Id, g[1].Type.String(),
					killedUnits, g[1].UnitCount, damage)
			}
			g[1].UnitCount -= killedUnits
			if g[1].UnitCount < 0 {
				g[1].UnitCount = 0
			}
		}
	}
	// Remove all killed groups.
	system.RemoveKilledGroups()
}

func part1() {
	system := loadScenario("input.txt")
	for len(system.Infection) > 0 && len(system.ImmuneSystem) > 0 {
		system.Print()
		FullAttack(system, true)
	}

	system.Print()
	sumGroup := system.ImmuneSystem
	if len(system.ImmuneSystem) == 0 {
		sumGroup = system.Infection
	}
	sum := 0
	for _, g := range sumGroup {
		sum += g.UnitCount
	}
	fmt.Println(sum)
}

// Return true if immune system is winning.
func Simulate(system *System) bool {
	for len(system.Infection) > 0 && len(system.ImmuneSystem) > 0 {
		FullAttack(system, false)
	}
	return len(system.ImmuneSystem) > 0
}

func part2() {
	// System hangs at 45 for some reason.
	system := &System{}
	var boost int
	for i := 50; i >= 46; i-- {
		system = loadScenario("input.txt")
		system.Boost(i)
		result := Simulate(system)
		fmt.Println(i, result)
		if !result {
			break
		}
		boost = i
	}
	fmt.Println("Boost:", boost)
	sum := 0
	for _, g := range system.ImmuneSystem {
		sum += g.UnitCount
	}
	fmt.Println("Unit count =", sum)
}

func main() {
	part1()
	part2()
}
