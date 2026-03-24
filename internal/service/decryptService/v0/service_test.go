package v0

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_DecryptValidService(t *testing.T) {
	data, err := os.ReadFile("testdata/private.pem")
	require.NoError(t, err)

	pemBlock, _ := pem.Decode(data)

	private, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	require.NoError(t, err)

	// Encrypt a test message using the public key
	testMessage := []byte("test message for decryption")
	encryptedMessage, err := rsa.EncryptOAEP(md5.New(), rand.Reader, &private.PublicKey, testMessage, nil)
	require.NoError(t, err)

	tests := []struct {
		name    string
		message []byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "valid encrypted message",
			message: encryptedMessage,
			want:    testMessage,
			wantErr: false,
		},
		{
			name:    "invalid encrypted message",
			message: []byte("invalid encrypted data"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty message",
			message: []byte{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "nil message",
			message: nil,
			want:    nil,
			wantErr: true,
		},
	}

	validService, err := New(Config{
		DecryptEnabled: true,
		CryptoKey:      "testdata/private.pem",
	})
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := validService.Decrypt(ctx, tt.message)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_DecryptDisabledService(t *testing.T) {
	tests := []struct {
		name    string
		message []byte
		want    []byte
	}{
		{
			name:    "returns message as-is",
			message: []byte("original message"),
			want:    []byte("original message"),
		},
	}

	disabledService, err := New(Config{
		DecryptEnabled: false,
		CryptoKey:      "testdata/private.pem",
	})
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := disabledService.Decrypt(ctx, tt.message)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_DecryptNoKeyService(t *testing.T) {
	tests := []struct {
		name    string
		message []byte
	}{
		{
			name:    "service without private key",
			message: []byte("any message"),
		},
	}

	noKeyService, err := New(Config{
		DecryptEnabled: true,
		CryptoKey:      "some_key.pem",
	})
	assert.Error(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			_, err = noKeyService.Decrypt(ctx, tt.message)

			assert.Error(t, err)
			assert.ErrorIs(t, err, errNotFound)
		})
	}
}
