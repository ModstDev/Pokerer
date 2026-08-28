package table

import (
	"context"
	"fmt"

	"github.com/ModstDev/Pokerer/internal/poker/game"
)

type Table struct {
	ID   string
	Game *game.Game

	actions chan ActionRequest
	done    chan struct{}
}

type ActionRequest struct {
	PlayerID string
	Action   game.Action
	Result   chan error
}

func NewTable(id string, g *game.Game) *Table {
	return &Table{
		ID:      id,
		Game:    g,
		actions: make(chan ActionRequest),
		done:    make(chan struct{}),
	}
}

func (t *Table) Run(ctx context.Context) {
	for {
		select {
		case request := <-t.actions:
			t.handleAction(request)

		case <-ctx.Done():
			return

		case <-t.done:
			return
		}
	}
}

func (t *Table) handleAction(request ActionRequest) {
	playerIndex := t.findPlayer(request.PlayerID)

	if playerIndex == -1 {
		request.Result <- fmt.Errorf("player is not at this table")
		return
	}

	if t.Game.CurrentPlayer != playerIndex {
		request.Result <- fmt.Errorf("it is not the player's turn")
		return
	}

	err := t.Game.ApplyAction(request.Action)

	request.Result <- err
}

func (t *Table) findPlayer(playerID string) int {
	for i := range t.Game.Players {
		if t.Game.Players[i].ID == playerID {
			return i
		}
	}

	return -1
}

func (t *Table) SubmitAction(ctx context.Context, request ActionRequest) error {
	request.Result = make(chan error, 1)

	select {
	case t.actions <- request:
	case <-ctx.Done():
		return ctx.Err()
	case <-t.done:
		return fmt.Errorf("table is closed")
	}

	select {
	case err := <-request.Result:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-t.done:
		return fmt.Errorf("table is closed")
	}
}
