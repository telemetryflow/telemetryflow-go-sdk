// =============================================================================
// Test Fixtures - Credentials
// =============================================================================
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.
//
// =============================================================================

package fixtures

// ValidCredentials contains valid test credentials
var ValidCredentials = struct {
	KeyID     string
	KeySecret string
}{
	KeyID:     "tfk_test_valid_key_id_12345",
	KeySecret: "tfs_test_valid_secret_67890abcdef",
}

// InvalidCredentials contains various invalid credential scenarios
var InvalidCredentials = []struct {
	Name      string
	KeyID     string
	KeySecret string
	Error     string
}{
	{
		Name:      "empty_key_id",
		KeyID:     "",
		KeySecret: "tfs_valid_secret",
		Error:     "key ID cannot be empty",
	},
	{
		Name:      "empty_key_secret",
		KeyID:     "tfk_valid_key",
		KeySecret: "",
		Error:     "key secret cannot be empty",
	},
	{
		Name:      "both_empty",
		KeyID:     "",
		KeySecret: "",
		Error:     "key ID cannot be empty",
	},
	{
		Name:      "invalid_key_id_prefix",
		KeyID:     "invalid_key_id",
		KeySecret: "tfs_valid_secret",
		Error:     "key ID must start with 'tfk_'",
	},
	{
		Name:      "invalid_key_secret_prefix",
		KeyID:     "tfk_valid_key",
		KeySecret: "invalid_secret",
		Error:     "key secret must start with 'tfs_'",
	},
	{
		Name:      "key_id_too_short",
		KeyID:     "tfk_",
		KeySecret: "tfs_valid_secret",
		Error:     "key ID is too short",
	},
	{
		Name:      "key_secret_too_short",
		KeyID:     "tfk_valid_key",
		KeySecret: "tfs_",
		Error:     "key secret is too short",
	},
	{
		Name:      "whitespace_key_id",
		KeyID:     "  tfk_valid_key  ",
		KeySecret: "tfs_valid_secret",
		Error:     "key ID contains invalid characters",
	},
	{
		Name:      "whitespace_key_secret",
		KeyID:     "tfk_valid_key",
		KeySecret: "  tfs_valid_secret  ",
		Error:     "key secret contains invalid characters",
	},
}

// CredentialFormats contains credentials in various formats for parsing tests
var CredentialFormats = struct {
	JSON string
	YAML string
	ENV  string
}{
	JSON: `{
  "key_id": "tfk_test_key_id",
  "key_secret": "tfs_test_secret"
}`,
	YAML: `key_id: tfk_test_key_id
key_secret: tfs_test_secret`,
	ENV: `TELEMETRYFLOW_API_KEY_ID=tfk_test_key_id
TELEMETRYFLOW_API_KEY_SECRET=tfs_test_secret`,
}
