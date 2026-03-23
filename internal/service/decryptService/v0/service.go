package v0

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

var errNotFound = errors.New("private key not found")

type Service struct {
	private *rsa.PrivateKey
	cfg     Config
}

func New(config Config) (*Service, error) {
	if !config.DecryptEnabled {
		return &Service{cfg: config}, nil
	}

	data, err := os.ReadFile(config.CryptoKey)
	if err != nil {
		return &Service{cfg: config}, fmt.Errorf("read private key error: %w", err)
	}

	pemBlock, _ := pem.Decode(data)
	if pemBlock == nil {
		return &Service{cfg: config}, errNotFound
	}

	private, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if err != nil {
		return &Service{cfg: config}, fmt.Errorf("parse private key error: %w", err)
	}

	return &Service{
		private: private,
		cfg:     config,
	}, nil
}

func (s *Service) Decrypt(ctx context.Context, message []byte) ([]byte, error) {
	if !s.cfg.DecryptEnabled {
		return message, nil
	}

	if s.private == nil {
		return nil, errNotFound
	}

	decrypted, err := rsa.DecryptOAEP(md5.New(), rand.Reader, s.private, message, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt message error: %w", err)
	}

	return decrypted, nil
}
