package game

import "fmt"

func (g *Game) bettingRoundComplete() bool {
	activePlayers := 0
	playersWhoCanAct := 0

	for _, player := range g.Players {
		if player.Folded {
			continue
		}

		activePlayers++

		if player.AllIn {
			continue
		}

		playersWhoCanAct++

		if !player.Acted {
			return false
		}

		if player.Bet != g.CurrentBet {
			return false
		}
	}

	if activePlayers <= 1 {
		return true
	}

	return playersWhoCanAct > 0
}

func (g *Game) resetActionsAfterRaise(raiser *Player) {
	for i := range g.Players {
		if &g.Players[i] == raiser {
			g.Players[i].Acted = true
			continue
		}

		if !g.Players[i].Folded && !g.Players[i].AllIn {
			g.Players[i].Acted = false
		}
	}
}

func (g *Game) setFirstPostFlopPlayer() error {
	for i := 1; i <= len(g.Players); i++ {
		index := nextPosition(
			g.DealerPosition,
			i,
			len(g.Players),
		)

		player := &g.Players[index]

		if !player.Folded && !player.AllIn {
			g.CurrentPlayer = index
			return nil
		}
	}

	return fmt.Errorf("no player can act")
}

func (g *Game) advancePlayer() error {
	if len(g.Players) == 0 {
		return fmt.Errorf("no players")
	}

	for i := 1; i <= len(g.Players); i++ {
		index := (g.CurrentPlayer + i) % len(g.Players)

		player := &g.Players[index]

		if !player.Folded && !player.AllIn {
			g.CurrentPlayer = index
			return nil
		}
	}

	return fmt.Errorf("no active players")
}

func (g *Game) resetBettingRound() {
	for i := range g.Players {
		g.Players[i].Bet = 0
		g.Players[i].Acted = false
	}

	g.CurrentBet = 0
}
