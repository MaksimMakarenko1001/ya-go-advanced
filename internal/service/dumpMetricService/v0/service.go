package v0

import (
	"fmt"
	"os"
	"sync"
)

type Service struct {
	mtx              sync.Mutex
	fname            string
	metricRepository MetricRepository
}

func New(fname string, metricRepo MetricRepository) *Service {
	return &Service{
		fname:            fname,
		metricRepository: metricRepo,
	}
}

func (srv *Service) ReadDump() error {
	file, err := os.ReadFile(srv.fname)
	if err != nil {
		return fmt.Errorf("reading file not ok, %w", err)
	}

	if err := srv.metricRepository.Load(file); err != nil {
		return fmt.Errorf("loading into repo not ok, %w", err)
	}

	return nil
}

func (srv *Service) WriteDump() error {
	srv.mtx.Lock()
	defer srv.mtx.Unlock()

	data, err := srv.metricRepository.Save()
	if err != nil {
		return fmt.Errorf("dump repo not ok, %w", err)
	}

	if err := os.WriteFile(srv.fname, data, 0666); err != nil {
		return fmt.Errorf("writting file not ok, %w", err)
	}

	return nil

}
