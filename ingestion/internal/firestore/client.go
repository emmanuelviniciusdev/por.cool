package firestore

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// ServiceAccount represents a Firebase service account JSON structure
type ServiceAccount struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

// Client handles communication with Firestore
type Client struct {
	serviceAccount *ServiceAccount
	httpClient     *http.Client
	accessToken    string
	tokenExpiry    time.Time
}

// SyncMetadataEntry represents an entry in the syncMetadata array
type SyncMetadataEntry struct {
	Name               string `json:"name"`
	LatestSyncDatetime string `json:"latestSyncDatetime"`
}

// NewClient creates a new Firestore client from a service account file
func NewClient(serviceAccountPath string) (*Client, error) {
	data, err := os.ReadFile(serviceAccountPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read service account file: %w", err)
	}

	var sa ServiceAccount
	if err := json.Unmarshal(data, &sa); err != nil {
		return nil, fmt.Errorf("failed to parse service account JSON: %w", err)
	}

	if sa.ProjectID == "" || sa.PrivateKey == "" || sa.ClientEmail == "" {
		return nil, fmt.Errorf("invalid service account: missing required fields")
	}

	return &Client{
		serviceAccount: &sa,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// NewClientFromServiceAccount creates a new Firestore client from a ServiceAccount struct
func NewClientFromServiceAccount(sa *ServiceAccount) (*Client, error) {
	if sa.ProjectID == "" || sa.PrivateKey == "" || sa.ClientEmail == "" {
		return nil, fmt.Errorf("invalid service account: missing required fields")
	}

	return &Client{
		serviceAccount: sa,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// base64URLEncode encodes bytes to base64 URL encoding without padding
func base64URLEncode(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}

// createJWT creates a signed JWT for Google OAuth2
func (c *Client) createJWT() (string, error) {
	header := map[string]string{
		"alg": "RS256",
		"typ": "JWT",
	}

	now := time.Now().Unix()
	claims := map[string]interface{}{
		"iss":   c.serviceAccount.ClientEmail,
		"sub":   c.serviceAccount.ClientEmail,
		"aud":   "https://oauth2.googleapis.com/token",
		"iat":   now,
		"exp":   now + 3600,
		"scope": "https://www.googleapis.com/auth/datastore",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	signatureInput := base64URLEncode(headerJSON) + "." + base64URLEncode(claimsJSON)

	// Parse the private key
	block, _ := pem.Decode([]byte(c.serviceAccount.PrivateKey))
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("private key is not RSA")
	}

	// Sign the JWT
	hash := sha256.Sum256([]byte(signatureInput))
	signature, err := rsa.SignPKCS1v15(nil, rsaKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signatureInput + "." + base64URLEncode(signature), nil
}

// getAccessToken retrieves an access token, refreshing if necessary
func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	// Return cached token if still valid
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	jwt, err := c.createJWT()
	if err != nil {
		return "", fmt.Errorf("failed to create JWT: %w", err)
	}

	data := url.Values{}
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	data.Set("assertion", jwt)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request access token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	// Set expiry to a bit before actual expiry to be safe
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second)

	return c.accessToken, nil
}

// listSettingsDocuments lists all documents in the settings collection
func (c *Client) listSettingsDocuments(ctx context.Context, accessToken string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://firestore.googleapis.com/v1/projects/%s/databases/(default)/documents/settings?pageSize=300",
		c.serviceAccount.ProjectID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list settings documents: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list settings failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Documents []map[string]interface{} `json:"documents"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Documents, nil
}

// parseFirestoreArrayPreservingAll parses the syncMetadata array preserving all entries
func parseFirestoreArrayPreservingAll(arrayValue map[string]interface{}) []SyncMetadataEntry {
	var result []SyncMetadataEntry

	values, ok := arrayValue["values"].([]interface{})
	if !ok {
		return result
	}

	for _, v := range values {
		mapValue, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		mv, ok := mapValue["mapValue"].(map[string]interface{})
		if !ok {
			continue
		}

		fields, ok := mv["fields"].(map[string]interface{})
		if !ok {
			continue
		}

		entry := SyncMetadataEntry{}

		if nameField, ok := fields["name"].(map[string]interface{}); ok {
			if nameStr, ok := nameField["stringValue"].(string); ok {
				entry.Name = nameStr
			}
		}

		if datetimeField, ok := fields["latestSyncDatetime"].(map[string]interface{}); ok {
			if datetimeStr, ok := datetimeField["timestampValue"].(string); ok {
				entry.LatestSyncDatetime = datetimeStr
			}
		}

		result = append(result, entry)
	}

	return result
}

// updateSettingsDocument updates the syncMetadata array in a settings document
func (c *Client) updateSettingsDocument(ctx context.Context, accessToken string, doc map[string]interface{}, syncServiceName string) error {
	currentDatetime := time.Now().UTC().Format(time.RFC3339)

	// Parse existing syncMetadata
	var syncMetadata []SyncMetadataEntry

	if fields, ok := doc["fields"].(map[string]interface{}); ok {
		if syncMetadataField, ok := fields["syncMetadata"].(map[string]interface{}); ok {
			if arrayValue, ok := syncMetadataField["arrayValue"].(map[string]interface{}); ok {
				syncMetadata = parseFirestoreArrayPreservingAll(arrayValue)
			}
		}
	}

	// Find existing entry for this sync service
	existingIndex := -1
	for i, entry := range syncMetadata {
		if entry.Name == syncServiceName {
			existingIndex = i
			break
		}
	}

	if existingIndex >= 0 {
		// Update existing entry
		syncMetadata[existingIndex].LatestSyncDatetime = currentDatetime
	} else {
		// Add new entry
		syncMetadata = append(syncMetadata, SyncMetadataEntry{
			Name:               syncServiceName,
			LatestSyncDatetime: currentDatetime,
		})
	}

	// Convert back to Firestore format
	values := make([]map[string]interface{}, len(syncMetadata))
	for i, entry := range syncMetadata {
		values[i] = map[string]interface{}{
			"mapValue": map[string]interface{}{
				"fields": map[string]interface{}{
					"name": map[string]interface{}{
						"stringValue": entry.Name,
					},
					"latestSyncDatetime": map[string]interface{}{
						"timestampValue": entry.LatestSyncDatetime,
					},
				},
			},
		}
	}

	syncMetadataFirestore := map[string]interface{}{
		"arrayValue": map[string]interface{}{
			"values": values,
		},
	}

	// Get document name
	docName, ok := doc["name"].(string)
	if !ok {
		return fmt.Errorf("document has no name field")
	}

	// Prepare batch write request
	batchWriteBody := map[string]interface{}{
		"writes": []map[string]interface{}{
			{
				"update": map[string]interface{}{
					"name": docName,
					"fields": map[string]interface{}{
						"syncMetadata": syncMetadataFirestore,
					},
				},
				"updateMask": map[string]interface{}{
					"fieldPaths": []string{"syncMetadata"},
				},
			},
		},
	}

	bodyJSON, err := json.Marshal(batchWriteBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("https://firestore.googleapis.com/v1/projects/%s/databases/(default)/documents:batchWrite",
		c.serviceAccount.ProjectID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// UpdateSettingsSyncMetadata updates the syncMetadata array in all settings documents
// This upserts an entry with the given name and current timestamp
func (c *Client) UpdateSettingsSyncMetadata(ctx context.Context, syncServiceName string) error {
	accessToken, err := c.getAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// List all settings documents
	docs, err := c.listSettingsDocuments(ctx, accessToken)
	if err != nil {
		return fmt.Errorf("failed to list settings documents: %w", err)
	}

	if len(docs) == 0 {
		return nil // No documents to update
	}

	// Update each settings document
	var lastErr error
	successCount := 0
	for _, doc := range docs {
		if err := c.updateSettingsDocument(ctx, accessToken, doc, syncServiceName); err != nil {
			lastErr = err
			continue
		}
		successCount++
	}

	if successCount == 0 && lastErr != nil {
		return fmt.Errorf("failed to update any settings documents: %w", lastErr)
	}

	return nil
}

// GetProjectID returns the project ID from the service account
func (c *Client) GetProjectID() string {
	return c.serviceAccount.ProjectID
}
