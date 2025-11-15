package v0

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/MaksimMakarenko1001/ya-go-advanced.git/pkg"
)

type Service struct {
	cfg Config
}

func New(config Config) *Service {
	return &Service{
		cfg: config,
	}
}

func (s *Service) Validate(ctx context.Context, message []byte, hash string) error {
	if s.cfg.Key == "" {
		return nil
	}

	hashBytes, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return pkg.ErrInternalServer.SetInfof("failed to decode hash, %v", err)
	}

	h := hmac.New(sha256.New, []byte(s.cfg.Key))
	if _, err := h.Write(message); err != nil {
		return pkg.ErrInternalServer.SetInfof("failed to validate message, %v", err)
	}

	if !hmac.Equal(h.Sum(nil), hashBytes) {
		return pkg.ErrBadRequest.SetInfo("invalid hash")
	}
	return nil
}

func (s *Service) Hash(ctx context.Context, message []byte) (string, error) {
	if s.cfg.Key == "" {
		return "", nil
	}

	h := hmac.New(sha256.New, []byte(s.cfg.Key))
	if _, err := h.Write(message); err != nil {
		return "", pkg.ErrInternalServer.SetInfof("failed to hash message, %v", err)
	}

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
