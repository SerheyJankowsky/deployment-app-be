# API Key Authentication System

This Go backend now includes a comprehensive API key authentication system that provides secure access for external integrations and programmatic access to your APIs.

## Features

✅ **Encrypted Storage**: API keys are encrypted in the database using AES-256-GCM encryption  
✅ **Secure Generation**: API keys are generated using cryptographically secure random number generation  
✅ **Multiple Auth Methods**: Support for both JWT and API key authentication  
✅ **Guard Middleware**: Easy-to-use middleware for protecting routes  
✅ **Management API**: Complete CRUD operations for API key management

## API Key Management

### Generate API Key

```bash
POST /api/v1/auth/generate-api-key
Authorization: Bearer <jwt_token>
```

**Response:**

```json
{
  "message": "API key generated successfully",
  "api_key": "64-character-hex-string",
  "user_id": 123
}
```

### Get API Key Status

```bash
GET /api/v1/auth/api-key
Authorization: Bearer <jwt_token>
```

**Response:**

```json
{
  "has_api_key": true,
  "api_key_preview": "abcd1234..."
}
```

### Revoke API Key

```bash
DELETE /api/v1/auth/revoke-api-key
Authorization: Bearer <jwt_token>
```

**Response:**

```json
{
  "message": "API key revoked successfully"
}
```

## Using API Keys

### Method 1: API-Key Header (Recommended)

```bash
curl -H "API-Key: your-64-character-api-key" \
     https://your-api.com/api/v1/api-secrets/
```

### Method 2: Authorization Bearer Header

```bash
curl -H "Authorization: Bearer your-64-character-api-key" \
     https://your-api.com/api/v1/api-secrets/
```

## Protected Routes

### API Key Only Routes

These routes only accept API key authentication:

- `GET /api/v1/api-secrets/` - List secrets
- `GET /api/v1/api-secrets/:id` - Get specific secret

### Combined Auth Routes

These routes accept both JWT and API key authentication:

- `POST /api/v1/api-secrets/` - Create secret

## Implementation Details

### Guard Middleware

#### API Key Guard

Protects routes with API key authentication only:

```go
apiKeyGuard := guards.ApiKeyGuard(userService)
router.Get("/protected", apiKeyGuard, handler)
```

#### Combined Guard

Accepts both JWT and API key authentication:

```go
combinedGuard := guards.CombinedGuard(userService)
router.Post("/flexible", combinedGuard, handler)
```

### Handler Context

In your handlers, you can access user information and authentication method:

```go
func MyHandler(ctx *fiber.Ctx) error {
    user := ctx.Locals("user")
    authMethod := ctx.Locals("auth_method") // "jwt" or "api_key"

    if authMethod == "api_key" {
        // Handle API key authentication
        if u, ok := user.(users.User); ok {
            // Access user.ID, user.IV, etc.
        }
    } else if authMethod == "jwt" {
        // Handle JWT authentication
        if claims, ok := user.(*libs.UserClaims); ok {
            // Access claims.UserID, claims.Email, etc.
        }
    }

    return ctx.Next()
}
```

## Security Features

### Encryption

- API keys are encrypted using AES-256-GCM encryption
- Each user has a unique IV (Initialization Vector)
- Encryption key is managed via environment variable `ENCRYPTION_KEY`

### Generation

- Uses `crypto/rand` for cryptographically secure random generation
- 32 bytes (256 bits) of entropy
- Encoded as 64-character hexadecimal string

### Database Storage

```sql
-- User table includes encrypted API key fields
api_key VARCHAR(255), -- Encrypted API key
iv VARCHAR(255)       -- Initialization vector for encryption
```

## Environment Variables

Ensure you have the following environment variable set:

```bash
ENCRYPTION_KEY=64-character-hex-string-for-aes-256-encryption
```

## Error Responses

### Missing API Key

```json
{
  "message": "API key required",
  "error": "Missing API-Key header"
}
```

### Invalid API Key

```json
{
  "message": "Invalid API key",
  "error": "API key not found or invalid"
}
```

### Combined Auth Error

```json
{
  "message": "Unauthorized",
  "error": "Valid JWT token or API key required"
}
```

## Best Practices

1. **Store API Keys Securely**: Never commit API keys to version control
2. **Use HTTPS**: Always use HTTPS in production to protect API keys in transit
3. **Rotate Keys**: Regularly regenerate API keys for security
4. **Limit Scope**: Use different API keys for different integrations
5. **Monitor Usage**: Track API key usage in your application logs

## Example Usage

### Python Example

```python
import requests

api_key = "your-64-character-api-key"
headers = {"API-Key": api_key}

response = requests.get(
    "https://your-api.com/api/v1/api-secrets/",
    headers=headers
)

if response.status_code == 200:
    secrets = response.json()
    print(f"Found {len(secrets)} secrets")
else:
    print(f"Error: {response.status_code} - {response.text}")
```

### JavaScript/Node.js Example

```javascript
const axios = require("axios");

const apiKey = "your-64-character-api-key";

async function getSecrets() {
  try {
    const response = await axios.get(
      "https://your-api.com/api/v1/api-secrets/",
      {
        headers: {
          "API-Key": apiKey,
        },
      }
    );

    console.log("Secrets:", response.data);
  } catch (error) {
    console.error("Error:", error.response?.data || error.message);
  }
}

getSecrets();
```

### cURL Examples

```bash
# Generate API key (requires JWT)
curl -X POST https://your-api.com/api/v1/auth/generate-api-key \
  -H "Authorization: Bearer your-jwt-token"

# Use API key to access protected resource
curl -H "API-Key: your-api-key" \
     https://your-api.com/api/v1/api-secrets/

# Create new secret with API key
curl -X POST https://your-api.com/api/v1/api-secrets/ \
  -H "API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{"name": "test", "content": "secret-value"}'
```

## Migration Notes

If you're adding this to an existing system:

1. The `User` model already includes `ApiKey` and `IV` fields
2. Existing users will have empty API keys initially
3. Users need to generate API keys through the management endpoints
4. The system is backward compatible with existing JWT authentication

This API key system provides a robust foundation for external integrations while maintaining security best practices.
