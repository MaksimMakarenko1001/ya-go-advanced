package remote

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

type Repository struct {
	url    string
	client *http.Client
}

func New(url string) *Repository {
	return &Repository{
		url:    url,
		client: &http.Client{},
	}
}

func (r *Repository) RemoteSend(ctx context.Context, body []byte) error {
	rq, err := http.NewRequestWithContext(ctx, http.MethodPost, r.url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("request not ok, %w", err)
	}

	rq.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(rq)
	if err != nil {
		return fmt.Errorf("http not ok, %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send, status_code=%d", resp.StatusCode)
	}
	return nil
}
