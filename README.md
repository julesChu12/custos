# User Domain

The **User Domain** is a core domain service that manages **user identity, lifecycle, security, and authorization**.  
It is distinct from the **Mora capability library** (which provides reusable technical modules) and the **API layer** (which orchestrates services and handles trust/zero trust).

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

### 4. Authorization & Access Control
- **RBAC** (Role-Based Access Control, via Casbin)
- **ABAC** (Attribute-Based Access Control, e.g., user can only access own data)
- Multi-tenant/organization support (SaaS B-end)
- APIs for managing user-role-permission relationships

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
- **Trust/Zero Trust** (device, IP, network validation → handled by API layer)
- **Infrastructure capabilities** (logging, config, DB wrappers → handled by Mora)
- **Other business domains** (orders, payments, etc.)

---

## Relationship with Mora & API Layer
- **Mora** → Provides reusable capabilities (auth token signing, logger, config, etc.)
- **User Domain** → Owns user model, lifecycle, security, roles/permissions
- **API Layer** → Orchestrates User Domain + Mora, enforces trust policies

---

## Roadmap
- ✅ MVP: User + Login + Simple Role  
- 🔒 Security: Forced logout, Refresh token rotation, Audit logs  
- 🔑 OAuth2.0: Third-party login, Provider support  
- 🏢 B-end: Casbin RBAC, multi-tenant, org model  