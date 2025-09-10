package keys

import (
	"crypto/ed25519"
	"crypto/rand"
	"novakeyauth/internal/keys"
	"testing"
)

func GenerateKey(t *testing.T) string{
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	pemPriv, err := keys.PrivateKeyToOpenSSHPEM(priv)
	if err != nil {
		t.Fatal(err)
	}	
	return pemPriv
}