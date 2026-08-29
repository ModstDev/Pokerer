package table

import (
	"context"
	"fmt"
	"sync"

	"github.com/ModstDev/Pokerer/internal/poker/game"
)

type Manager struct {
	mu     sync.RWMutex
	tables map[string]*Table
	ctx    context.Context
	loader GameLoader
}

type GameLoader interface {
	LoadGame(ctx context.Context, tableID string) (*game.Game, error)
}

func NewManager(ctx context.Context, loader GameLoader) *Manager {
	return &Manager{
		tables: make(map[string]*Table),
		ctx:    ctx,
		loader: loader,
	}
}

func (m *Manager) Add(table *Table) error {
	if table == nil {
		return fmt.Errorf("table is nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.tables[table.ID]; exists {
		return fmt.Errorf("table already exists")
	}

	m.tables[table.ID] = table

	go table.Run(m.ctx)

	return nil
}

func (m *Manager) Get(id string) (*Table, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	table, ok := m.tables[id]

	return table, ok
}

func (m *Manager) Remove(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	table, exists := m.tables[id]
	if !exists {
		return false
	}

	table.Close()

	delete(m.tables, id)

	return true
}

func (m *Manager) Create(id string, g *game.Game) (*Table, error) {
	table := NewTable(id, g)

	if err := m.Add(table); err != nil {
		return nil, err
	}

	return table, nil
}

func (m *Manager) GetOrCreate(ctx context.Context, id string) (*Table, error) {
	table, ok := m.Get(id)
	if ok {
		return table, nil
	}

	g, err := m.loader.LoadGame(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("loading game: %w", err)
	}

	table = NewTable(id, g)

	m.mu.Lock()
	defer m.mu.Unlock()

	// Another goroutine may have created the table
	// while we were loading the game.
	if existing, ok := m.tables[id]; ok {
		return existing, nil
	}

	m.tables[id] = table

	go table.Run(m.ctx)

	return table, nil
}
