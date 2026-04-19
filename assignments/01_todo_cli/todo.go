package todo

// Todo is the main application struct that holds dependencies.
type Todo struct {
	Storer Storer
}

// Option is a functional option for configuring Todo.
type Option func(*Todo)

// WithFileStorer sets a FileStore with the given path.
func WithFileStorer(path string) Option {
	return func(t *Todo) {
		t.Storer = NewFileStore(path)
	}
}

// WithMemoryStorer sets an in-memory store.
func WithMemoryStorer() Option {
	return func(t *Todo) {
		t.Storer = NewMemoryStore()
	}
}

// New creates a new Todo with sensible defaults.
// Default store uses "tasks.json".
func New(opts ...Option) *Todo {
	t := &Todo{
		Storer: NewMemoryStore(),
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}
