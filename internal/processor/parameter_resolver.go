package processor

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/b0nbon1/stratal/internal/security"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type ParameterResolver struct {
	store         *db.SQLStore
	secretManager *security.SecretManager
}

func NewParameterResolver(store *db.SQLStore, secretManager *security.SecretManager) *ParameterResolver {
	return &ParameterResolver{
		store:         store,
		secretManager: secretManager,
	}
}

// ResolveParameters resolves task parameters, including secrets and task outputs
func (pr *ParameterResolver) ResolveParameters(
	ctx context.Context,
	task db.Task,
	userID pgtype.UUID,
	taskOutputs map[string]string,
) (map[string]string, map[string]string, error) {

	resolvedParams := make(map[string]string)
	secretEnvVars := make(map[string]string)

	// 1. Copy regular parameters and resolve ${TASK_OUTPUT.task_name} references
	for key, value := range task.Config.Parameters {
		resolvedValue := pr.resolveTaskOutputReferences(value, taskOutputs)
		resolvedParams[key] = resolvedValue
	}

	// 2. Resolve secrets
	for secretName, envVarName := range task.Config.Secrets {
		secret, err := pr.store.GetSecretByName(ctx, db.GetSecretByNameParams{
			Name:   secretName,
			UserID: userID,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("secret '%s' not found: %w", secretName, err)
		}

		decrypted, err := pr.secretManager.Decrypt(secret.EncryptedValue)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decrypt secret '%s': %w", secretName, err)
		}

		secretEnvVars[envVarName] = decrypted
	}

	return resolvedParams, secretEnvVars, nil
}

// resolveTaskOutputReferences replaces ${TASK_OUTPUT.task_name} with actual values
func (pr *ParameterResolver) resolveTaskOutputReferences(value string, taskOutputs map[string]string) string {
	re := regexp.MustCompile(`\$\{TASK_OUTPUT\.([^}]+)\}`)
	return re.ReplaceAllStringFunc(value, func(match string) string {
		taskName := strings.TrimPrefix(strings.TrimSuffix(match, "}"), "${TASK_OUTPUT.")
		if output, exists := taskOutputs[taskName]; exists {
			return output
		}
		return match // Keep original if not found
	})
}
