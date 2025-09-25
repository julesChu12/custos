# User Domain (Custos)

The **User Domain (Custos)** is a core domain service that manages **user identity, lifecycle, security, and authorization**.  
It is distinct from the **Mora capability library** (generic modules like auth/logger/config) and the **API layer** (Clotho, which orchestrates services and enforces trust).

---

## Responsibilities

### 1. User Lifecycle Management
- User registration (C-end self-service, B-end admin-created)
- Activation / Freeze / Deletion
- Profile management (nickname, avatar, email, phone, extended profile)

### 2. Identity Authentication
- Username + password login
- Phone/email OTP login (C-end)
- OAuth2.0 third-party login (Google, WeChat, Apple ID, etc.)
- Access/Refresh token mechanism (supports rotation + state table)
- Multi-session management (web, mobile, tablet)
- Forced logout (based on `token_version` or session state)

### 3. User Security
- Password hashing (bcrypt/argon2)
- Login failure limit (anti-brute-force)
- Two-factor authentication (2FA/MFA: SMS, email, OTP)
- Login & audit logs
- Abnormal login detection (geo-based, concurrent sessions, repeated failures)

### 4. Authorization & Access Control (via Casbin)
- RBAC implemented with **Casbin**
- Custos does **not** maintain custom `roles/permissions` tables
- Instead:
  - `users` table stores user accounts
  - Casbin `casbin_rule` table stores role & permission policies
  - Custos integrates Casbin Enforcer for runtime checks
  - Custos provides wrapper APIs for managing roles/permissions by calling Casbin APIs
- Future: ABAC (Attribute-Based Access Control) using Casbin models

### 5. OAuth2.0 Provider Capabilities (Optional)
- Provide `/authorize`, `/token`, `/userinfo`
- Support Grant Types: Authorization Code, Client Credentials, Refresh Token
- For open platform / external API integrations

### 6. Audit & Observability
- Login events (success/failure, IP, UA)
- Permission change logs
- Security events (forced logout, reused refresh token detection)
- Export to MQ/ES/Prometheus for auditing and alerting

---

## Out of Scope
The User Domain **does not handle**:
- **Trust/Zero Trust** (device, IP, network validation → handled by API layer Clotho)
- **Infrastructure capabilities** (logging, config, DB wrappers → handled by Mora)
- **Other business domains** (orders, payments, etc.)

---

## Database Schema (DDL)

### users
```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(64) UNIQUE NOT NULL,
    email VARCHAR(128) UNIQUE,
    phone VARCHAR(32) UNIQUE,
    password_hash VARCHAR(255),
    user_type ENUM('customer','staff','partner') DEFAULT 'customer',
    tenant_id BIGINT NULL,
    status ENUM('active','disabled','locked','deleted') DEFAULT 'active',
    token_version INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### user_profiles
```sql
CREATE TABLE user_profiles (
    user_id BIGINT PRIMARY KEY,
    nickname VARCHAR(64),
    avatar VARCHAR(255),
    gender ENUM('male','female','other') DEFAULT 'other',
    birthday DATE,
    extra JSON,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### user_oauth
```sql
CREATE TABLE user_oauth (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    provider VARCHAR(64) NOT NULL,
    provider_uid VARCHAR(128) NOT NULL,
    access_token VARCHAR(255),
    refresh_token VARCHAR(255),
    expires_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider, provider_uid),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### refresh_tokens
```sql
CREATE TABLE refresh_tokens (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    token_hash CHAR(64) NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### casbin_rule (used by Casbin Adapter)
```sql
CREATE TABLE casbin_rule (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    ptype VARCHAR(32) NOT NULL,
    v0 VARCHAR(255),
    v1 VARCHAR(255),
    v2 VARCHAR(255),
    v3 VARCHAR(255),
    v4 VARCHAR(255),
    v5 VARCHAR(255)
);
```

---

## Roadmap
- ✅ MVP: users table + login + refresh token rotation + Casbin RBAC integration
- 🔒 Security: Forced logout, audit logs, MFA
- 🔑 OAuth2.0: Third-party login, Provider support
- 🏢 B-end: Multi-tenant, org model
