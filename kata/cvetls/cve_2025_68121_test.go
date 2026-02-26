// Package cvetls holds unit tests that validate CVE fixes in Go's crypto/tls.
//
// CVE-2025-68121: Config.Clone must not copy automatically generated session
// ticket keys, so that two Configs do not share keys and session resumption
// cannot occur across them when it should be isolated.
// See: https://go.dev/issue/77113
package cvetls

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"
)

// generateTestCert creates a self-signed certificate for server TLS tests.
func generateTestCert(t *testing.T, commonName string) tls.Certificate {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		t.Fatalf("serial: %v", err)
	}
	template := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: commonName},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{commonName},
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	return tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  key,
	}
}

// TestCVE_2025_68121_CloneDoesNotShareAutoSessionTicketKeys verifies that
// after Config.Clone(), the cloned config does not share automatically
// generated session ticket keys with the original. So a session ticket
// issued by server A cannot be used to resume a session with server B
// that was created by cloning A (when neither A nor B set SessionTicketKey
// or SetSessionTicketKeys).
//
// With the CVE fixed (Go 1.25.7+): resumption to B must not happen (DidResume false).
// With the CVE present: B would accept the ticket and DidResume could be true.
func TestCVE_2025_68121_CloneDoesNotShareAutoSessionTicketKeys(t *testing.T) {
	const serverName = "test.cve.example.com"
	cert := generateTestCert(t, serverName)

	// Server config A: no explicit SessionTicketKey or SetSessionTicketKeys,
	// so crypto/tls will auto-generate session ticket keys when used.
	configA := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		// SessionTicketKey left unset; do not call SetSessionTicketKeys.
	}

	// Listen with config A and run one handshake so that auto-generated
	// session ticket keys exist.
	lnA, err := tls.Listen("tcp", "127.0.0.1:0", configA)
	if err != nil {
		t.Fatalf("listen A: %v", err)
	}
	addrA := lnA.Addr().String()
	defer lnA.Close()

	// Client session cache shared across both connections.
	sessionCache := tls.NewLRUClientSessionCache(10)
	clientConfig := &tls.Config{
		ServerName:         serverName,
		InsecureSkipVerify: true,
		ClientSessionCache: sessionCache,
		MinVersion:         tls.VersionTLS12,
	}

	// First connection: client connects to A, full handshake, caches session.
	go func() {
		conn, _ := lnA.Accept()
		if conn != nil {
			_ = conn.(*tls.Conn).Handshake()
			_ = conn.Close()
		}
	}()
	connToA, err := tls.Dial("tcp", addrA, clientConfig)
	if err != nil {
		t.Fatalf("dial A: %v", err)
	}
	if err := connToA.Handshake(); err != nil {
		t.Fatalf("handshake A: %v", err)
	}
	// First connection is never a resumption.
	if connToA.ConnectionState().DidResume {
		t.Fatal("first connection to A should not be resumption")
	}
	_ = connToA.Close()

	// Clone A to B (after A has been used, so auto keys exist).
	configB := configA.Clone()

	// Server B with cloned config.
	lnB, err := tls.Listen("tcp", "127.0.0.1:0", configB)
	if err != nil {
		t.Fatalf("listen B: %v", err)
	}
	addrB := lnB.Addr().String()
	defer lnB.Close()

	// Second connection: same client (same session cache) connects to B.
	// With CVE fixed: B does not have A's auto-generated keys, so the cached
	// ticket is not valid for B; client must do full handshake (DidResume false).
	// With CVE present: B would have copied keys, resumption could succeed (DidResume true).
	go func() {
		conn, _ := lnB.Accept()
		if conn != nil {
			_ = conn.(*tls.Conn).Handshake()
			_ = conn.Close()
		}
	}()
	connToB, err := tls.Dial("tcp", addrB, clientConfig)
	if err != nil {
		t.Fatalf("dial B: %v", err)
	}
	if err := connToB.Handshake(); err != nil {
		t.Fatalf("handshake B: %v", err)
	}
	_ = connToB.Close()

	if connToB.ConnectionState().DidResume {
		t.Error("CVE-2025-68121: session was resumed with cloned config B; Clone() must not copy auto-generated session ticket keys (expected full handshake, DidResume=false)")
	}
}
