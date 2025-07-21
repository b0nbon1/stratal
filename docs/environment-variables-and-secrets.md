# Environment Variables and Secrets Management

This document describes how to configure and use environment variables and secrets in Stratal.

## Environment Variables

### Infrastructure Configuration

Stratal uses environment variables for infrastructure configuration. Create a `.env` file in the root directory:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=stratal_user
DB_PASSWORD=secure_db_password
DB_NAME=stratal
DB_SSL_MODE=disable

# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=secure_redis_password
REDIS_DB=0

# Security
ENCRYPTION_KEY=your-32-character-encryption-key!!

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
```

### Required Environment Variables

- `ENCRYPTION_KEY`: 32-character key for encrypting user secrets (required)

### Optional Environment Variables

All other variables have sensible defaults for development.

## User Secrets Management

### Creating Secrets

Use the API to create encrypted secrets for users:

```bash
curl -X POST http://localhost:8080/api/v1/secrets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "smtp_password",
    "value": "your-secret-smtp-password"
  }'
```

### Listing Secrets

```bash
curl http://localhost:8080/api/v1/secrets
```

### Using Secrets in Tasks

Reference secrets in task configurations:

```json
{
  "name": "send_email",
  "type": "builtin",
  "config": {
    "parameters": {
      "smtp_host": "smtp.gmail.com",
      "smtp_port": "587",
      "from": "noreply@company.com",
      "to": "user@example.com",
      "subject": "Test Email"
    },
    "secrets": {
      "smtp_password": "SMTP_PASSWORD"
    }
  }
}
```

## Task Parameter Resolution

### Parameter Types

1. **Regular Parameters**: Static values or interpolated task outputs
2. **Secrets**: Encrypted values stored in the database
3. **Task Outputs**: Results from previous tasks

### Parameter Interpolation

Use `${TASK_OUTPUT.task_name}` to reference outputs from previous tasks:

```json
{
  "parameters": {
    "email": "${TASK_OUTPUT.get_user_email}",
    "message": "Hello ${TASK_OUTPUT.get_user_name}!"
  }
}
```

### Environment Variables in Custom Scripts

Custom scripts receive:

1. **Regular parameters**: As environment variables
2. **Secrets**: As environment variables (names defined in config)
3. **Task outputs**: As `TASK_OUTPUT_TASK_NAME` environment variables

```python
import os

# Access regular parameter
api_endpoint = os.environ.get('API_ENDPOINT')

# Access secret
api_key = os.environ.get('API_KEY')

# Access previous task output
user_id = os.environ.get('TASK_OUTPUT_GET_USER_ID')
```

## Security Features

### Encryption

- All user secrets are encrypted using AES-256-GCM
- Encryption key must be 32 characters long
- Each secret has a unique nonce for security

### Access Control

- Secrets are scoped to users (user_id)
- Only the owning user can access their secrets
- Secrets are decrypted only during task execution

### Audit Trail

- All secret operations are logged
- Creation and access timestamps are tracked
- Failed decryption attempts are logged

## Best Practices

### Environment Variables

1. **Never commit `.env` files** to version control
2. **Use different encryption keys** for different environments
3. **Rotate encryption keys** periodically
4. **Use strong, random passwords** for database and Redis

### Secrets

1. **Use descriptive secret names** (e.g., "smtp_gmail_password")
2. **Rotate secrets regularly**
3. **Use least privilege principle** for secret access
4. **Monitor secret usage** through logs

### Task Configuration

1. **Separate secrets from parameters** in task configs
2. **Use meaningful environment variable names** for secrets
3. **Document required secrets** in job descriptions
4. **Test with dummy secrets** in development

## Deployment

### Development

1. Copy `.env.example` to `.env`
2. Generate a 32-character encryption key
3. Configure local database and Redis

### Production

1. **Use a secrets management service** (AWS Secrets Manager, etc.)
2. **Set environment variables** through your deployment system
3. **Use strong encryption keys** (generated, not hardcoded)
4. **Enable TLS** for database and Redis connections
5. **Monitor and rotate** encryption keys

## Troubleshooting

### Common Issues

1. **"ENCRYPTION_KEY must be exactly 32 characters"**
   - Ensure your encryption key is exactly 32 characters long

2. **"Secret not found"**
   - Verify the secret name matches exactly
   - Check that the secret belongs to the correct user

3. **"Failed to decrypt secret"**
   - The encryption key may have changed
   - The secret data may be corrupted

4. **Environment variables not loaded**
   - Ensure `.env` file is in the root directory
   - Check file permissions (readable by the process)

### Debugging

Enable debug logging to trace secret resolution:

```bash
export LOG_LEVEL=debug
./stratal
```
```

This implementation provides:

1. **Comprehensive Configuration Management**: Using `godotenv` for environment variables
2. **Secure Secret Storage**: AES-256-GCM encryption for user secrets
3. **Parameter Resolution**: Automatic interpolation of secrets and task outputs
4. **API Management**: Full CRUD operations for secrets
5. **Enhanced Task Execution**: Support for secrets in both builtin and custom tasks
6. **Security Best Practices**: Proper encryption, access control, and audit trails

To use this system:

1. Run `go mod tidy` to install dependencies
2. Run `sqlc generate` to generate the new secret queries
3. Create a `.env` file with your configuration
4. Start the server and worker with the new secret management capabilities

The system maintains backward compatibility while adding powerful secret management features for secure job execution. 