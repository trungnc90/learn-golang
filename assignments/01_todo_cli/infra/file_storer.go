package infra

import (
	"encoding/json"
	"os"

	todo "github.com/trungnc90/learn-golang/assignments/01_todo_cli"
)

// FileStore implements Storer using a JSON file.
type FileStore struct {
	FilePath string
}

func NewFileStore(path string) *FileStore {
	return &FileStore{FilePath: path}
}

func (fs *FileStore) Load() ([]todo.Task, error) {
	data, err := os.ReadFile(fs.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []todo.Task{}, nil
		}
		return nil, err
	}

	var tasks []todo.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (fs *FileStore) Save(tasks []todo.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fs.FilePath, data, 0644)
}
