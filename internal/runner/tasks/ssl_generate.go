package tasks

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caddyserver/certmagic"
)

// SSLGenerateTask generates SSL certificates using CertMagic and stores them in a temporary folder
func SSLGenerateTask(ctx context.Context, params map[string]string) (string, error) {
	// Required parameters
	domain, exists := params["domain"]
	if !exists || domain == "" {
		return "", fmt.Errorf("missing required parameter: domain")
	}

	email, exists := params["email"]
	if !exists || email == "" {
		return "", fmt.Errorf("missing required parameter: email")
	}

	// Optional parameters with defaults
	storageDir := params["storage_dir"]
	if storageDir == "" {
		// Create a temporary directory for certificate storage
		tempDir, err := os.MkdirTemp("", "ssl-certs-*")
		if err != nil {
			return "", fmt.Errorf("failed to create temporary directory: %w", err)
		}
		storageDir = tempDir
	}

	// Optional: staging mode for testing (Let's Encrypt staging)
	staging := params["staging"] == "true"

	// Optional: key type (default: "rsa2048")
	keyType := params["key_type"]
	if keyType == "" {
		keyType = "rsa2048"
	}

	// Check context cancellation before proceeding
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}


	// Configure key type
	var keyTypeValue certmagic.KeyType
	switch strings.ToLower(keyType) {
	case "rsa2048", "rsa":
		keyTypeValue = certmagic.RSA2048
	case "rsa4096":
		keyTypeValue = certmagic.RSA4096
	case "p256", "ec256", "ecdsa":
		keyTypeValue = certmagic.P256
	case "p384", "ec384":
		keyTypeValue = certmagic.P384
	case "ed25519":
		keyTypeValue = certmagic.ED25519
	default:
		return "", fmt.Errorf("unsupported key type: %s (supported: rsa2048, rsa4096, p256, p384, ed25519)", keyType)
	}
	
	// Create a custom config with our key generator
	config := certmagic.NewDefault()
	config.Storage = &certmagic.FileStorage{Path: storageDir}
	config.KeySource = &certmagic.StandardKeyGenerator{
		KeyType: keyTypeValue,
	}
	
	// Create ACME issuer
	acmeTemplate := certmagic.ACMEIssuer{
		Email:  email,
		Agreed: true,
	}
	
	if staging {
		acmeTemplate.CA = certmagic.LetsEncryptStagingCA
	}
	
	acmeIssuer := certmagic.NewACMEIssuer(config, acmeTemplate)
	config.Issuers = []certmagic.Issuer{acmeIssuer}

	// Split domains if multiple are provided (comma-separated)
	domains := strings.Split(domain, ",")
	for i, d := range domains {
		domains[i] = strings.TrimSpace(d)
	}

	// Check context cancellation before certificate generation
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Obtain certificates
	err := config.ManageSync(ctx, domains)
	if err != nil {
		return "", fmt.Errorf("failed to obtain SSL certificates: %w", err)
	}

	// Verify certificates were created and get their paths
	var certPaths []string
	caEndpoint := acmeIssuer.CA
	if caEndpoint == "" {
		caEndpoint = certmagic.LetsEncryptProductionCA
	}
	
	for _, d := range domains {
		// CertMagic stores certificates in a specific structure
		certPath := filepath.Join(storageDir, "certificates", caEndpoint, d, d+".crt")
		keyPath := filepath.Join(storageDir, "certificates", caEndpoint, d, d+".key")
		
		// Check if files exist, but don't fail if they're in a different location
		// CertMagic might use a different internal structure
		var actualCertPath, actualKeyPath string
		
		// Try to find the actual certificate files
		err := filepath.Walk(filepath.Join(storageDir, "certificates"), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Continue walking even if some paths fail
			}
			
			if strings.Contains(path, d) {
				if strings.HasSuffix(path, ".crt") || strings.HasSuffix(path, ".pem") {
					actualCertPath = path
				} else if strings.HasSuffix(path, ".key") {
					actualKeyPath = path
				}
			}
			return nil
		})
		
		if err != nil {
			// If we can't walk the directory, use the expected paths
			actualCertPath = certPath
			actualKeyPath = keyPath
		}
		
		// If we couldn't find the files through walking, use expected paths
		if actualCertPath == "" {
			actualCertPath = certPath
		}
		if actualKeyPath == "" {
			actualKeyPath = keyPath
		}
		
		certPaths = append(certPaths, fmt.Sprintf("Certificate: %s\nPrivate Key: %s", actualCertPath, actualKeyPath))
	}

	// Build result message
	var result strings.Builder
	result.WriteString(fmt.Sprintf("SSL certificates generated successfully for domains: %s\n", strings.Join(domains, ", ")))
	result.WriteString(fmt.Sprintf("Storage directory: %s\n", storageDir))
	result.WriteString(fmt.Sprintf("Key type: %s\n", keyType))
	if staging {
		result.WriteString("Mode: Staging (Let's Encrypt Staging CA)\n")
	} else {
		result.WriteString("Mode: Production (Let's Encrypt Production CA)\n")
	}
	result.WriteString("\nGenerated files:\n")
	for _, path := range certPaths {
		result.WriteString(fmt.Sprintf("  %s\n", path))
	}

	return result.String(), nil
}