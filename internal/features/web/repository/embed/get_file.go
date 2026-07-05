package web_embed_repository

import (
	"errors"
	"fmt"
	"io/fs"

	core_errors "github.com/vladislav-koval/todo-app/internal/core/errors"
)

func (r *WebRepository) GetFile(name string) ([]byte, error) {
	file, err := files.ReadFile(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("file %q: %w", name, core_errors.ErrNotFound)
		}

		return nil, fmt.Errorf("file %q: %w", name, err)
	}

	return file, nil
}
