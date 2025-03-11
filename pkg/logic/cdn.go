package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/anthrove/identity/pkg/object"
	"path/filepath"
	"strings"
)

// ServeResource retrieves a specific resource within a specified tenant.
//
// Parameters:
//   - ctx: context for managing request-scoped values, cancelation, and deadlines.
//   - tenantID: unique identifier of the tenant to which the resource belongs.
//   - filePath: unique identifier of the resource to be retrieved.
//
// Returns:
//   - The local resource path.
//   - Error if there is any issue during retrieval.
func (is IdentityService) ServeResource(ctx context.Context, tenantID string, resourcePath string) (string, error) {
	providers, err := is.FindProviders(ctx, tenantID, object.Pagination{})
	if err != nil {
		return "", err
	}
	sanitizedFilePath := sanitizeFilePath(resourcePath)

	var parameters map[string]string

	for _, provider := range providers {
		if provider.Category == "storage" && provider.ProviderType == "local" {
			err := json.Unmarshal(provider.Parameter, &parameters)
			if err != nil {
				return "", err
			}
		}

		sanitizedBasePath := sanitizeFilePath(parameters["base_path"])
		localFilePath := fmt.Sprintf("local_storage_provider/%s/%s/%s", tenantID, sanitizedBasePath, sanitizedFilePath)

		return localFilePath, nil
	}

	return "", fmt.Errorf("no suitable provider found")
}

func sanitizeFilePath(path string) string {
	return strings.TrimPrefix(filepath.Clean(path), "/")
}
