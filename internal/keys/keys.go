package keys

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
)

func EncodeSSHPublicKey(pub ed25519.PublicKey, username string) (string, error) {
	if len(pub) != ed25519.PublicKeySize {
		return "", fmt.Errorf("invalid public key size: %d", len(pub))
	}

	var buf bytes.Buffer

	if err := writeString(&buf, "ssh-ed25519"); err != nil {
		return "", err
	}

	if err := writeBytes(&buf, pub); err != nil {
		return "", err
	}

	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	if username == "" {
		return fmt.Sprintf("ssh-ed25519 %s", b64), nil
	}
	return fmt.Sprintf("ssh-ed25519 %s %s", b64, username), nil
}

func writeString(buf *bytes.Buffer, s string) error {
	if err := binary.Write(buf, binary.BigEndian, uint32(len(s))); err != nil {
		return err
	}
	if _, err := buf.Write([]byte(s)); err != nil {
		return err
	}
	return nil
}

func writeBytes(buf *bytes.Buffer, b []byte) error {
	if err := binary.Write(buf, binary.BigEndian, uint32(len(b))); err != nil {
		return err
	}
	if _, err := buf.Write(b); err != nil {
		return err
	}
	return nil
}

func ParseOpenSSHED25519Signer(keyPEM string) (ssh.Signer, ed25519.PublicKey, error) {
	signer, err := ssh.ParsePrivateKey([]byte(keyPEM))
	if err != nil {
		return nil, nil, fmt.Errorf("parse openssh key: %w", err)
	}

	if signer.PublicKey().Type() != "ssh-ed25519" {
		return nil, nil, errors.New("not an ed25519 key")
	}

	// Extract the public key
	if ed25519PubKey, ok := signer.PublicKey().(ssh.CryptoPublicKey); ok {
		if cryptoPub, ok := ed25519PubKey.CryptoPublicKey().(ed25519.PublicKey); ok {
			return signer, cryptoPub, nil
		}
	}

	return nil, nil, errors.New("could not extract ed25519 public key")
}

func PrivateKeyToOpenSSHPEM(priv ed25519.PrivateKey) (string, error) {
	if len(priv) != ed25519.PrivateKeySize {
		return "", fmt.Errorf("invalid ed25519 private key size: %d", len(priv))
	}

	pub := priv.Public().(ed25519.PublicKey)

	// Generate random checkint for integrity check
	checkint := make([]byte, 4)
	rand.Read(checkint)

	// Helper function to write length-prefixed strings/data
	writeBytes := func(buf *bytes.Buffer, data []byte) {
		binary.Write(buf, binary.BigEndian, uint32(len(data)))
		buf.Write(data)
	}

	writeString := func(buf *bytes.Buffer, s string) {
		writeBytes(buf, []byte(s))
	}

	// Build the unencrypted private key section
	privateSection := &bytes.Buffer{}
	
	// Write check integers (must be the same for unencrypted keys)
	privateSection.Write(checkint)
	privateSection.Write(checkint)
	
	// Write the key data
	writeString(privateSection, "ssh-ed25519")
	writeBytes(privateSection, pub)
	
	// For ED25519, the private key is stored as seed(32) + public(32) = 64 bytes
	privateKeyBlob := make([]byte, 64)
	copy(privateKeyBlob[:32], priv.Seed())
	copy(privateKeyBlob[32:], pub)
	writeBytes(privateSection, privateKeyBlob)
	
	// Empty comment
	writeString(privateSection, "")
	
	// Padding to make total length multiple of cipher block size (8 for unencrypted)
	padLen := (8 - (privateSection.Len() % 8)) % 8
	for i := 1; i <= padLen; i++ {
		privateSection.WriteByte(byte(i))
	}

	// Build the public key section
	publicSection := &bytes.Buffer{}
	writeString(publicSection, "ssh-ed25519")
	writeBytes(publicSection, pub)

	// Build the main structure
	result := &bytes.Buffer{}
	
	// Magic bytes
	result.WriteString("openssh-key-v1\x00")
	
	// Cipher, KDF, and options (all empty for unencrypted)
	writeString(result, "none")         // cipher
	writeString(result, "none")         // kdf
	writeBytes(result, []byte{})        // kdf options (empty)
	
	// Number of keys
	binary.Write(result, binary.BigEndian, uint32(1))
	
	// Public key section
	writeBytes(result, publicSection.Bytes())
	
	// Private key section
	writeBytes(result, privateSection.Bytes())

	// Base64 encode the result
	encoded := base64.StdEncoding.EncodeToString(result.Bytes())

	// Format as PEM with proper line wrapping
	var pemResult bytes.Buffer
	pemResult.WriteString("-----BEGIN OPENSSH PRIVATE KEY-----\n")
	
	for len(encoded) > 70 {
		pemResult.WriteString(encoded[:70] + "\n")
		encoded = encoded[70:]
	}
	if len(encoded) > 0 {
		pemResult.WriteString(encoded + "\n")
	}
	
	pemResult.WriteString("-----END OPENSSH PRIVATE KEY-----\n")

	return pemResult.String(), nil
}