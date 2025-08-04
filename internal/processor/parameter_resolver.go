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

// resolveTaskOutputReferences replaces ${TASK_OUTPUT.task_name} and ${task_name.output} with actual values
func (pr *ParameterResolver) resolveTaskOutputReferences(value string, taskOutputs map[string]string) string {
	// Handle ${TASK_OUTPUT.task_name} pattern
	re1 := regexp.MustCompile(`\$\{TASK_OUTPUT\.([^}]+)\}`)
	value = re1.ReplaceAllStringFunc(value, func(match string) string {
		taskName := strings.TrimPrefix(strings.TrimSuffix(match, "}"), "${TASK_OUTPUT.")
		if output, exists := taskOutputs[taskName]; exists {
			return output
		}
		return match // Keep original if not found
	})

	// Handle ${task_name.output} pattern
	re2 := regexp.MustCompile(`\$\{([^}]+)\.output\}`)
	value = re2.ReplaceAllStringFunc(value, func(match string) string {
		taskName := strings.TrimSuffix(strings.TrimPrefix(match, "${"), ".output}")
		if output, exists := taskOutputs[taskName]; exists {
			return output
		}
		return match // Keep original if not found
	})

	return value
}
