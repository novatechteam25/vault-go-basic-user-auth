# Vault Go Basic Auth

Simple authentication implementation using HashiCorp Vault with Go and Gin framework.

## 📋 Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Setup in Vault](#setup-in-vault)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)
- [Architecture](#architecture)

## 🎯 Overview

This project demonstrates:
- **Token-based authentication** using Vault
- **Role-based access control** (Admin vs Regular users)
- **No database dependency** - Vault is the single source of truth
- **Protected routes** with middleware validation
- **Policy-driven authorization**

## 📦 Prerequisites

- Go 1.19+
- Vault server running at `https://dev-vault.fromnovatech.xyz/`
- Valid Vault root or admin token

## 🚀 Installation

```bash
# Clone the repository
git clone <repo-url>
cd vault-go-basicauth

# Download dependencies
go mod download

# Run the application
go run main.go
```

## ⚙️ Setup in Vault

### 1. Enable Token Auth (usually default)

```bash
vault auth enable token
```

### 2. Create Policies in Vault UI

Navigate to `Policies > ACL Policies` and create:

#### `my-app-policy`
```hcl
path "secret/data/my-app/*" {
  capabilities = ["read", "list"]
}

path "auth/token/renew-self" {
  capabilities = ["update"]
}

path "auth/token/lookup-self" {
  capabilities = ["read"]
}
```

#### `admin-policy` (optional)
```hcl
path "*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
```

#### `regular-policy` (optional)
```hcl
path "secret/data/user/*" {
  capabilities = ["read", "list"]
}

path "auth/token/renew-self" {
  capabilities = ["update"]
}

path "auth/token/lookup-self" {
  capabilities = ["read"]
}
```

### 3. Or Create via CLI

```bash
# Login first
vault login <your-token>

# Create policy
vault policy write my-app-policy - <<EOF
path "secret/data/my-app/*" {
  capabilities = ["read", "list"]
}

path "auth/token/renew-self" {
  capabilities = ["update"]
}

path "auth/token/lookup-self" {
  capabilities = ["read"]
}
EOF
```

## 📡 API Endpoints

| Endpoint | Method | Auth Required | Description |
|----------|--------|---------------|-------------|
| `/public` | GET | No | Public route (no auth needed) |
| `/admin-only` | GET | Yes (admin) | Admin-only access |
| `/user-dashboard` | GET | Yes (regular) | Regular user access |
| `/my-protected-route` | GET | Yes | Protected route (any valid token) |

## 🧪 Testing

### 1. Start the application

```bash
go run main.go
```

You'll see output:
```
New token created: hvs.XXXXXXXXXXXXXXXXXXXXXX
Token Lease Duration: 3600 seconds
```

### 2. Test endpoints

**Public route (no auth):**
```bash
curl http://localhost:8080/public
```

**Protected route with token:**
```bash
TOKEN=hvs.XXXXXXXXXXXXXXXXXXXXXX
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/my-protected-route
```

**Admin-only route:**
```bash
curl -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:8080/admin-only
```

**Regular user route:**
```bash
curl -H "Authorization: Bearer $REGULAR_TOKEN" http://localhost:8080/user-dashboard
```

## 📊 Architecture

### Authentication Flow

```
CLIENT REQUEST
    ↓
    └─→ Authorization Header Check
        ├─→ Missing? → 401 Unauthorized
        └─→ Found?
            ↓
            Extract & Clean Token
            ↓
            Send to Vault (LookupSelf)
            ├─→ Invalid/Expired? → 401 Unauthorized
            └─→ Valid?
                ├─→ Check Policies (if needed)
                │   ├─→ Has required policy? → Continue
                │   └─→ No policy? → 403 Forbidden
                └─→ Execute Route Handler
                    ↓
                    Return Response (200)
```

### Token Structure

```
Token (hvs.XXXXXXXXXXXXXXXXXXXXXX)
├── Policies: ["admin-policy", "regular-policy"]
├── Display Name: "admin-token" or "regular-token"
├── TTL: 1 hour (renewable)
└── Created by: Server on startup
```

## 📝 Response Examples

**Success (200):**
```json
{
  "message": "Access granted to protected route!"
}
```

**Missing Authorization (401):**
```json
{
  "error": "Missing Authorization header"
}
```

**Invalid Token (401):**
```json
{
  "error": "Invalid token"
}
```

**Insufficient Permissions (403):**
```json
{
  "error": "Insufficient permissions for this route"
}
```

## 🔐 Key Features

- ✅ **No Database Needed** - Vault manages all auth
- ✅ **Token Expiration** - Configurable TTL (default 1h)
- ✅ **Role-Based Access** - Different permissions for different users
- ✅ **Token Renewal** - Tokens can be renewed before expiration
- ✅ **Middleware Pattern** - Reusable auth middleware
- ✅ **Stateless** - Perfect for microservices

## 📚 Token Lifecycle

1. **Creation** - App creates token on startup with specific policies
2. **Usage** - Client includes token in Authorization header
3. **Validation** - Middleware validates token with Vault
4. **Expiration** - After TTL, token becomes invalid (default 1h)
5. **Renewal** - Token can be renewed if renewable flag is true

## 🤝 Contributing

Feel free to modify and extend this for your needs!

## 📄 License

MIT