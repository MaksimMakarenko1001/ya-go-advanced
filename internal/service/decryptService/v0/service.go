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

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/service/decryptService/service"
)

var _ service.DecryptService = (*Service)(nil)

var (
	errNotFound   = errors.New("not found")
	errPrivateKey = errors.New("private key error")
)

type Service struct {
	private     *rsa.PrivateKey
	privateFunc func(string) (*rsa.PrivateKey, error)
	cfg         Config
}

func New(config Config) *Service {
	s := &Service{cfg: config}

	s.privateFunc = func(cryptoKey string) (*rsa.PrivateKey, error) {
		data, err := os.ReadFile(cryptoKey)
		if err != nil {
			return nil, fmt.Errorf("read error: %w", err)
		}

		pemBlock, _ := pem.Decode(data)
		if pemBlock == nil {
			return nil, errNotFound
		}

		private, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse error: %w", err)
		}

		return private, nil
	}

	return s
}

func (s *Service) Decrypt(ctx context.Context, message []byte) ([]byte, error) {
	if !s.cfg.DecryptEnabled {
		return message, nil
	}

	if s.private == nil {
		private, err := s.privateFunc(s.cfg.CryptoKey)
		if err != nil {
			return nil, fmt.Errorf("%w:%s", errPrivateKey, err.Error())
		}

		s.private = private
	}

	decrypted, err := rsa.DecryptOAEP(md5.New(), rand.Reader, s.private, message, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt message error: %w", err)
	}

	return decrypted, nil
}
