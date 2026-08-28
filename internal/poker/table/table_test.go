package table

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/ModstDev/Pokerer/internal/poker/game"
)

func TestSubmitAction(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g := game.NewGame(game.GameConfig{
		SmallBlind: 10,
		BigBlind:   20,
	})

	if err := g.AddPlayer(game.Player{
		ID:    "alice",
		Seat:  0,
		Chips: 1000,
	}); err != nil {
		t.Fatal(err)
	}

	if err := g.AddPlayer(game.Player{
		ID:    "bob",
		Seat:  1,
		Chips: 1000,
	}); err != nil {
		t.Fatal(err)
	}

	rng := rand.New(rand.NewPCG(42, 0))

	if err := g.StartHand(rng); err != nil {
		t.Fatal(err)
	}

	table := NewTable("table-1", g)

	go table.Run(ctx)

	err := table.SubmitAction(
		ctx,
		ActionRequest{
			PlayerID: "alice",
			Action: game.Action{
				Type: game.ActionFold,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if !g.Players[0].Folded {
		t.Fatal("expected Alice to be folded")
	}

	if g.State != game.StateShowdown {
		t.Fatalf(
			"expected showdown, got %s",
			g.State,
		)
	}
}

func TestCloseTable(t *testing.T) {
	g := game.NewGame(game.GameConfig{
		SmallBlind: 10,
		BigBlind:   20,
	})

	table := NewTable("table-1", g)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})

	go func() {
		table.Run(ctx)
		close(done)
	}()

	table.Close()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("table did not stop after close")
	}
}

func TestCloseTableIsIdempotent(t *testing.T) {
	g := game.NewGame(game.GameConfig{
		SmallBlind: 10,
		BigBlind:   20,
	})

	table := NewTable("table-1", g)

	table.Close()
	table.Close()
}

func TestSubmitActionClosedTable(t *testing.T) {
	g := game.NewGame(game.GameConfig{
		SmallBlind: 10,
		BigBlind:   20,
	})

	table := NewTable("table-1", g)

	table.Close()

	ctx := context.Background()

	err := table.SubmitAction(
		ctx,
		ActionRequest{
			PlayerID: "alice",
			Action: game.Action{
				Type: game.ActionFold,
			},
		},
	)

	if err == nil {
		t.Fatal("expected error when submitting to closed table")
	}
}
