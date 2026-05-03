package model

import (
	"time"

	"gorm.io/gorm"
)

// HostconfJwk hold public and private JWKs
type HostconfJwk struct {
	gorm.Model
	KeyId        string    // JWK KID
	ExpiresAt    time.Time // Expiration time stamp
	PublicJwk    string    // Public JWK as serialized JSON
	EncryptionId string    // id of the encryption key
	EncryptedJwk []byte    // Encrypted private key, nil if key is revoked
}