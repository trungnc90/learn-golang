package todo

import (
	"encoding/json"
	"os"
)

// FileStore implements Store using a JSON file.
type FileStore struct {
	FilePath string
}

func NewFileStore(path string) *FileStore {
	return &FileStore{FilePath: path}
}

func (fs *FileStore) Load() ([]Task, error) {
	data, err := os.ReadFile(fs.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (fs *FileStore) Save(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fs.FilePath, data, 0644)
}
