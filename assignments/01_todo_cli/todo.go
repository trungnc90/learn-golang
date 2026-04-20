package todo

// Todo is the main application struct that holds dependencies.
type Todo struct {
	store Storer
}

// Option is a functional option for configuring Todo.
type Option func(*Todo)

// WithStorer sets a custom Storer implementation.
func WithStorer(s Storer) Option {
	return func(t *Todo) {
		t.store = s
	}
}

// New creates a new Todo with sensible defaults.
// Default store uses in-memory.
func New(opts ...Option) *Todo {
	t := &Todo{
		store: NewMemoryStore(),
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}
