package uuid

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
)

type base58Encoder struct{}

func (enc base58Encoder) Encode(u uuid.UUID) string {
	return base58.Encode(u[:])
}

func (enc base58Encoder) Decode(s string) (uuid.UUID, error) {
	return uuid.FromBytes(base58.Decode(s))
}

var enc58 base58Encoder
