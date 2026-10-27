package task

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var ErrNotFound = errors.New("task not found")

type Repo struct {
	mu    sync.RWMutex
	seq   int64
	items map[int64]*Task
}

func NewRepo() *Repo {
	return &Repo{items: make(map[int64]*Task)}
}

func (r *Repo) List() []*Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*Task, 0, len(r.items))
	for _, t := range r.items {
		copied := *t
		out = append(out, &copied)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (r *Repo) Get(id int64) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	copied := *t
	return &copied, nil
}

func (r *Repo) Create(title string) *Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.seq++
	now := time.Now().UTC()
	t := &Task{
		ID:        r.seq,
		Title:     strings.TrimSpace(title),
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.items[t.ID] = t

	copied := *t
	return &copied
}

func (r *Repo) Update(id int64, title *string, done *bool) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	if title != nil {
		t.Title = strings.TrimSpace(*title)
	}
	if done != nil {
		t.Done = *done
	}
	t.UpdatedAt = time.Now().UTC()

	copied := *t
	return &copied, nil
}

func (r *Repo) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	return nil
}
