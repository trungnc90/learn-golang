package todo

// MemoryStore implements Store using an in-memory slice.
// Used for testing — no file I/O, no cleanup needed.
type MemoryStore struct {
	tasks []Task
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{tasks: []Task{}}
}

func (ms *MemoryStore) Load() ([]Task, error) {
	// Return a copy to avoid unintended mutations
	result := make([]Task, len(ms.tasks))
	copy(result, ms.tasks)
	return result, nil
}

func (ms *MemoryStore) Save(tasks []Task) error {
	ms.tasks = make([]Task, len(tasks))
	copy(ms.tasks, tasks)
	return nil
}
