package file

import (
	"context"
	"errors"
	"os"
)

var errNoFileOpen = errors.New("error no file open")

type Repository struct {
	fname string
	file  *os.File
}

func New(fname string) *Repository {
	return &Repository{fname: fname}
}

func (r *Repository) FileOpen(ctx context.Context) error {
	file, err := os.OpenFile(r.fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return err
	}

	r.file = file
	return nil
}

func (r *Repository) FileAppend(ctx context.Context, line []byte) error {
	if r.file == nil {
		return errNoFileOpen
	}

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
