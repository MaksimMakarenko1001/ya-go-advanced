package service

import "context"

type DecryptService interface {
	Decrypt(ctx context.Context, message []byte) ([]byte, error)
}
