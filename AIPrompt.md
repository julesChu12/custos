You are an experienced software engineer and product manager. 
I am building a Golang monorepo with a capability library called **Mora** 
and a domain service called **User Domain**. 

Mora provides generic technical modules (auth, logger, config, db, cache, mq, utils) 
and adapters (gin/go-zero middleware), plus starter demos. 
It does NOT own business logic.

The **User Domain** is a domain service with these responsibilities:
1. User lifecycle (register, activate, freeze, delete, profile management)
2. Authentication (username+password, OTP login, OAuth2.0 third-party login, 
   access/refresh token with rotation, multi-session management, forced logout)
3. Security (password hashing, login failure limits, 2FA/MFA, audit logs, abnormal login detection)
4. Authorization (RBAC via Casbin, optional ABAC, user-role-permission relationships, multi-tenant)
5. Optional OAuth2.0 Provider endpoints (/authorize, /token, /userinfo)
6. Audit & Observability (login logs, permission changes, security events, export to MQ/ES/Prometheus)

The User Domain does NOT handle:
- Trust/Zero Trust (device, IP, network checks) → belongs to API Layer
- Infrastructure capabilities → belong to Mora
- Other business domains (orders, payments, etc.)

Tasks you might generate:
- Data model design (SQL schemas, Golang structs)
- Service APIs (Gin/Go-Zero handlers)
- Middleware (auth, RBAC checks)
- Refresh token rotation with state table
- Forced logout with token_version strategy
- Integration with Casbin for RBAC
- Integration with Mora for auth token signing/validation
- Sample starter code for login, logout, token refresh, and permission checks

Follow clean architecture principles: 
- Mora = capability library
- User Domain = business domain
- API Layer = orchestrator (glue)

Now, based on this context, generate code, architecture, and documentation.