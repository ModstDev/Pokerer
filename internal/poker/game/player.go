package game

type Player struct {
	ID                string
	Seat              int
	Chips             int64
	Bet               int64
	Cards             []Card
	Folded            bool
	AllIn             bool
	Acted             bool
	TotalContribution int64
}

func (g *Game) contribute(player *Player, amount int64) {
	player.Chips -= amount
	player.Bet += amount
	player.TotalContribution += amount
	g.Pot += amount
}
