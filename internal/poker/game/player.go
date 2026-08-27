package game

type Player struct {
	ID     string
	Seat   int
	Chips  int64
	Bet    int64
	Cards  []Card
	Folded bool
	AllIn  bool
	Acted  bool
}
