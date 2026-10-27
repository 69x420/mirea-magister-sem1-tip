package storage

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrNotFound = errors.New("task not found")

type Task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MemoryStore struct {
	mu    sync.RWMutex
	auto  int64
	tasks map[int64]*Task
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{tasks: make(map[int64]*Task)}
}

func (s *MemoryStore) Create(title string) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.auto++
	now := time.Now().UTC()
	t := &Task{
		ID:        s.auto,
		Title:     title,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.tasks[t.ID] = t
	return t
}

func (s *MemoryStore) Get(id int64) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	copied := *t
	return &copied, nil
}

func (s *MemoryStore) List(q string) []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	q = strings.ToLower(strings.TrimSpace(q))
	out := make([]*Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		if q != "" && !strings.Contains(strings.ToLower(t.Title), q) {
			continue
		}
		copied := *t
		out = append(out, &copied)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (s *MemoryStore) Update(id int64, title *string, done *bool) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	if title != nil {
		t.Title = *title
	}
	if done != nil {
		t.Done = *done
	}
	t.UpdatedAt = time.Now().UTC()

	copied := *t
	return &copied, nil
}

func (s *MemoryStore) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.tasks, id)
	return nil
}
