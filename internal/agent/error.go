package agent

import (
	"errors"
	"net/url"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/pkg/backoff"
)

func ClassifyHTTPError(err error) backoff.ErrorClassification {
	if err == nil {
		return backoff.NonRetriable
	}

	var urlError *url.Error
	if errors.As(err, &urlError) {
		return backoff.Retriable
	}

	return backoff.NonRetriable
}
