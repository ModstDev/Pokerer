package game

import "sort"

type Pot struct {
	Amount          int64
	EligiblePlayers []int
}

func (g *Game) BuildPots() []Pot {
	levels := make([]int64, 0)

	for _, player := range g.Players {
		if player.TotalContribution <= 0 {
			continue
		}

		levels = append(levels, player.TotalContribution)
	}

	sort.Slice(levels, func(i, j int) bool {
		return levels[i] < levels[j]
	})

	levels = uniqueInt64(levels)

	pots := make([]Pot, 0)

	var previousLevel int64

	for _, level := range levels {
		contribution := level - previousLevel

		if contribution <= 0 {
			continue
		}

		eligible := make([]int, 0)

		for i, player := range g.Players {
			if player.TotalContribution >= level && !player.Folded {
				eligible = append(eligible, i)
			}
		}

		amount := contribution * int64(len(eligible))

		if amount > 0 && len(eligible) >= 2 {
			pots = append(pots, Pot{
				Amount:          amount,
				EligiblePlayers: eligible,
			})
		}

		previousLevel = level
	}

	return pots
}

func uniqueInt64(values []int64) []int64 {
	if len(values) == 0 {
		return nil
	}

	result := []int64{values[0]}

	for i := 1; i < len(values); i++ {
		if values[i] != values[i-1] {
			result = append(result, values[i])
		}
	}

	return result
}
