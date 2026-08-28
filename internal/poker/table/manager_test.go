package table

import (
	"context"
	"testing"

	"github.com/ModstDev/Pokerer/internal/poker/game"
)

func TestManagerAddAndGet(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager := NewManager(ctx)

	g := game.NewGame(game.GameConfig{
		SmallBlind: 10,
		BigBlind:   20,
	})

	table := NewTable("table-1", g)

	if err := manager.Add(table); err != nil {
		t.Fatal(err)
	}

	result, ok := manager.Get("table-1")

	if !ok {
		t.Fatal("expected table to exist")
	}

	if result != table {
		t.Fatal("expected returned table to be the same table")
	}
}

func TestManagerRejectsDuplicateTable(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager := NewManager(ctx)

	g1 := game.NewGame(game.GameConfig{
		SmallBlind: 10,
		BigBlind:   20,
	})

	g2 := game.NewGame(game.GameConfig{
		SmallBlind: 10,
		BigBlind:   20,
	})

	if _, err := manager.Create("table-1", g1); err != nil {
		t.Fatal(err)
	}

	if _, err := manager.Create("table-1", g2); err == nil {
		t.Fatal("expected duplicate table error")
	}
}

func TestManagerRemove(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager := NewManager(ctx)

	g := game.NewGame(game.GameConfig{
		SmallBlind: 10,
		BigBlind:   20,
	})

	if _, err := manager.Create("table-1", g); err != nil {
		t.Fatal(err)
	}

	if !manager.Remove("table-1") {
		t.Fatal("expected table to be removed")
	}

	if _, ok := manager.Get("table-1"); ok {
		t.Fatal("expected table to no longer exist")
	}

	if manager.Remove("table-1") {
		t.Fatal("expected second removal to return false")
	}
}
