package file

import (
	"context"
	"errors"
	"os"
	"sync"
)

var errNoFileOpen = errors.New("error no file open")

type Repository struct {
	mx   *sync.Mutex
	file *os.File
}

func New(fname string) *Repository {
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return &Repository{}
	}

	return &Repository{
		mx:   &sync.Mutex{},
		file: file,
	}
}

func (r *Repository) FileAppend(ctx context.Context, line []byte) error {
	if r.file == nil {
		return errNoFileOpen
	}

	r.mx.Lock()
	defer r.mx.Unlock()

	line = append(line, '\n')
	_, err := r.file.Write(line)
	return err
}

func (r *Repository) FileClose(ctx context.Context) error {
	if r.file == nil {
		return nil
	}

	if err := r.file.Sync(); err != nil {
		return err
	}

	return r.file.Close()
}
