package protocol

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// TrustLevel represents how much a certificate is trusted
type TrustLevel string

const (
	// TrustPermanent means the certificate is permanently trusted
	TrustPermanent TrustLevel = "permanent"

	// TrustSession means the certificate is trusted for this session only
	TrustSession TrustLevel = "session"

	// TrustOnce means the certificate was accepted once but should prompt again
	TrustOnce TrustLevel = "once"
)

// CertificateInfo stores information about a known certificate
type CertificateInfo struct {
	// Fingerprint is the SHA256 fingerprint of the certificate
	Fingerprint string `json:"fingerprint"`

	// FirstSeen is when the certificate was first seen
	FirstSeen time.Time `json:"first_seen"`

	// LastSeen is when the certificate was last seen
	LastSeen time.Time `json:"last_seen"`

	// Trust indicates the trust level
	Trust TrustLevel `json:"trust"`

	// NotAfter is the certificate expiration date
	NotAfter time.Time `json:"not_after"`

	// Subject is the certificate subject
	Subject string `json:"subject"`
}

// KnownHosts stores the known certificates for each host
type KnownHosts struct {
	Version string                      `json:"version"`
	Hosts   map[string]*CertificateInfo `json:"hosts"`
}

// TOFUVerifier implements Trust On First Use certificate verification
type TOFUVerifier struct {
	mu         sync.RWMutex
	knownHosts *KnownHosts
	filePath   string

	// OnCertificateChange is called when a certificate changes
	// It should return true to accept the new certificate, false to reject
	OnCertificateChange func(hostname string, old, new *CertificateInfo) (bool, TrustLevel)

	// OnFirstSeen is called when a certificate is seen for the first time
	// It should return true to accept the certificate, false to reject
	OnFirstSeen func(hostname string, info *CertificateInfo) (bool, TrustLevel)
}

// NewTOFUVerifier creates a new TOFU verifier
func NewTOFUVerifier(filePath string) (*TOFUVerifier, error) {
	verifier := &TOFUVerifier{
		filePath: filePath,
		knownHosts: &KnownHosts{
			Version: "1.0",
			Hosts:   make(map[string]*CertificateInfo),
		},
	}

	// Try to load existing known hosts
	if err := verifier.Load(); err != nil {
		// If the file doesn't exist, that's okay
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load known hosts: %w", err)
		}
	}

	return verifier, nil
}

// VerifyCertificate verifies a certificate for a given hostname
func (v *TOFUVerifier) VerifyCertificate(hostname string, state tls.ConnectionState) error {
	if len(state.PeerCertificates) == 0 {
		return fmt.Errorf("no peer certificates")
	}

	cert := state.PeerCertificates[0]
	fingerprint := certificateFingerprint(cert)

	info := &CertificateInfo{
		Fingerprint: fingerprint,
		FirstSeen:   time.Now(),
		LastSeen:    time.Now(),
		Trust:       TrustPermanent,
		NotAfter:    cert.NotAfter,
		Subject:     cert.Subject.String(),
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	known, exists := v.knownHosts.Hosts[hostname]

	if !exists {
		// First time seeing this host
		accept, trustLevel := true, TrustPermanent

		if v.OnFirstSeen != nil {
			accept, trustLevel = v.OnFirstSeen(hostname, info)
		}

		if !accept {
			return fmt.Errorf("certificate rejected by user")
		}

		info.Trust = trustLevel
		v.knownHosts.Hosts[hostname] = info

		// Save to disk
		if err := v.save(); err != nil {
			return fmt.Errorf("failed to save known hosts: %w", err)
		}

		return nil
	}

	// Check if certificate has changed
	if known.Fingerprint != fingerprint {
		// Certificate has changed!
		accept, trustLevel := false, TrustOnce

		if v.OnCertificateChange != nil {
			accept, trustLevel = v.OnCertificateChange(hostname, known, info)
		}

		if !accept {
			return fmt.Errorf("certificate changed and was rejected")
		}

		info.FirstSeen = known.FirstSeen // Preserve first seen time
		info.Trust = trustLevel
		v.knownHosts.Hosts[hostname] = info

		// Save to disk
		if err := v.save(); err != nil {
			return fmt.Errorf("failed to save known hosts: %w", err)
		}

		return nil
	}

	// Certificate matches, update last seen time
	known.LastSeen = time.Now()

	// Save to disk (we could optimize this to not save on every request)
	if err := v.save(); err != nil {
		// Don't fail the request if we can't save
		// Just log it or handle it somehow
		_ = err
	}

	return nil
}

// certificateFingerprint computes the SHA256 fingerprint of a certificate
func certificateFingerprint(cert *x509.Certificate) string {
	hash := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(hash[:])
}

// Load loads the known hosts from disk
func (v *TOFUVerifier) Load() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	data, err := os.ReadFile(v.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v.knownHosts)
}

// save saves the known hosts to disk (caller must hold lock)
func (v *TOFUVerifier) save() error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(v.filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(v.knownHosts, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(v.filePath, data, 0600)
}

// Save saves the known hosts to disk (thread-safe version)
func (v *TOFUVerifier) Save() error {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.save()
}

// GetCertificateInfo returns information about a known certificate
func (v *TOFUVerifier) GetCertificateInfo(hostname string) (*CertificateInfo, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	info, exists := v.knownHosts.Hosts[hostname]
	return info, exists
}

// RemoveCertificate removes a certificate from the known hosts
func (v *TOFUVerifier) RemoveCertificate(hostname string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	delete(v.knownHosts.Hosts, hostname)
	return v.save()
}

// ClearAll removes all known certificates
func (v *TOFUVerifier) ClearAll() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.knownHosts.Hosts = make(map[string]*CertificateInfo)
	return v.save()
}
