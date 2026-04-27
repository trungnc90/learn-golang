package todo

// Todo implements the Manager interface.
type Todo struct {
	store Storer
}

// NewManager creates a new Todo with the given Storer.
func NewManager(store Storer) *Todo {
	return &Todo{store: store}
}
