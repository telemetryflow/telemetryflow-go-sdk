<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-dark.svg">
    <source media="(prefers-color-scheme: light)" srcset="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-light.svg">
    <img src="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-light.svg" alt="TelemetryFlow Logo" width="80%">
  </picture>

  <h3>TelemetryFlow GO SDK</h3>

[![Version](https://img.shields.io/badge/Version-1.1.1-orange.svg)](CHANGELOG.md)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/)
[![OTEL SDK](https://img.shields.io/badge/OpenTelemetry_SDK-1.39.0-blueviolet)](https://opentelemetry.io/)
[![OpenTelemetry](https://img.shields.io/badge/OTLP-100%25%20Compliant-success?logo=opentelemetry)](https://opentelemetry.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://hub.docker.com/r/telemetryflow/telemetryflow-sdk)

</div>

<p align="center">
  Enterprise-grade Go SDK for <a href="https://telemetryflow.id">TelemetryFlow</a> - the observability platform that provides unified metrics, logs, and traces collection following OpenTelemetry standards.
</p>

---

# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.1.x   | :white_check_mark: |
| 1.0.x   | :x:                |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of TelemetryFlow Core seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Where to Report

**Please DO NOT report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to:
- **Security Team**: security@telemetryflow.id
- **Project Lead**: support@telemetryflow.id

### What to Include

Please include the following information in your report:

- **Type of vulnerability** (e.g., SQL injection, XSS, authentication bypass)
- **Full paths of source file(s)** related to the vulnerability
- **Location of the affected source code** (tag/branch/commit or direct URL)
- **Step-by-step instructions** to reproduce the issue
- **Proof-of-concept or exploit code** (if possible)
- **Impact of the vulnerability** and how an attacker might exploit it

### Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Fix Timeline**: Depends on severity
  - Critical: 1-7 days
  - High: 7-14 days
  - Medium: 14-30 days
  - Low: 30-90 days

### Disclosure Policy

- Security issues will be disclosed after a fix is available
- We will credit researchers who report vulnerabilities (unless they prefer to remain anonymous)
- We follow responsible disclosure practices

## Security Best Practices

### For Users

#### 1. Environment Variables
```bash
# Never commit .env files
echo ".env" >> .gitignore

# Use strong secrets
pnpm run generate:secrets
```

#### 2. Database Security
```bash
# Use strong passwords
POSTGRES_PASSWORD=<strong-random-password>
CLICKHOUSE_PASSWORD=<strong-random-password>

# Restrict database access
# Only allow connections from trusted IPs
```

#### 3. JWT Configuration
```bash
# Use minimum 32 characters for secrets
JWT_SECRET=<min-32-chars-random-string>
SESSION_SECRET=<min-32-chars-random-string>

# Set appropriate expiration
JWT_EXPIRES_IN=24h  # Adjust based on your needs
```

#### 4. Production Deployment
```bash
# Always use NODE_ENV=production
NODE_ENV=production

# Disable debug logs
LOG_LEVEL=warn

# Enable HTTPS only
# Use reverse proxy (nginx/traefik) with SSL/TLS
```

### For Contributors

#### 1. Code Security

**Never commit:**
- Passwords or API keys
- Private keys or certificates
- Database credentials
- JWT secrets
- Personal information

**Always:**
- Use environment variables for sensitive data
- Validate all user inputs
- Sanitize database queries
- Use parameterized queries (TypeORM handles this)
- Implement proper authentication and authorization

#### 2. Dependencies

```bash
# Check for vulnerabilities
pnpm audit

# Fix vulnerabilities
pnpm audit fix

# Update dependencies regularly
pnpm update
```

#### 3. Code Review

All code changes must:
- Pass security review
- Include tests for security-critical features
- Follow OWASP security guidelines
- Be reviewed by at least one maintainer

## Security Features

### Authentication & Authorization

- **JWT-based authentication** with secure token generation
- **5-tier RBAC system** (Super Admin, Admin, Developer, Viewer, Demo)
- **Permission-based access control** with 22+ granular permissions
- **Password hashing** using Argon2 (industry standard)
- **Session management** with secure session secrets

### Data Protection

- **PostgreSQL** for transactional data with row-level security
- **ClickHouse** for audit logs and observability data
- **Encrypted connections** between services
- **Input validation** using class-validator
- **SQL injection prevention** via TypeORM parameterized queries

### Observability & Monitoring

- **Audit logging** for all critical operations
- **OpenTelemetry tracing** for request tracking
- **Winston logging** with structured logs
- **Health checks** for service monitoring

### Network Security

- **Docker network isolation** (172.151.151.0/24)
- **Service-to-service communication** on private network
- **Exposed ports** only for necessary services
- **CORS configuration** for API access control

## Vulnerability Disclosure

### Past Vulnerabilities

No security vulnerabilities have been reported yet.

### Security Advisories

Security advisories will be published at:
- GitHub Security Advisories
- Project documentation
- Release notes

## Compliance

### Standards

TelemetryFlow Core follows:
- **OWASP Top 10** security guidelines
- **CWE/SANS Top 25** vulnerability prevention
- **NIST Cybersecurity Framework** principles

### Certifications

Currently pursuing:
- SOC 2 Type II compliance
- ISO 27001 certification

## Security Contacts

### Primary Contact
- **Email**: security@telemetryflow.id
- **Response Time**: 48 hours

### Alternative Contact
- **Email**: support@telemetryflow.id
- **GitHub**: [@telemetryflow](https://github.com/telemetryflow)

## Bug Bounty Program

We currently do not have a formal bug bounty program, but we:
- Acknowledge security researchers in release notes
- Provide public recognition for valid reports
- Consider monetary rewards for critical vulnerabilities (case-by-case basis)

## Security Updates

### Notification Channels

Stay informed about security updates:
- **GitHub Releases**: Watch repository for releases
- **Security Advisories**: Enable GitHub security alerts
- **Changelog**: Check [CHANGELOG.md](./CHANGELOG.md)
- **Release Notes**: Review [docs/RELEASE_NOTES_*.md](./docs/)

### Update Process

```bash
# Check current version
cat package.json | grep version

# Update to latest version
git pull origin main
pnpm install

# Run migrations if needed
pnpm db:migrate

# Restart services
docker-compose restart
```

## Contribution Security Guidelines

### Before Contributing

1. **Read** [CONTRIBUTING.md](./CONTRIBUTING.md)
2. **Review** this security policy
3. **Sign** commits with GPG key (recommended)
4. **Test** security implications of your changes

### Code Submission

```bash
# Sign commits
git commit -S -m "Your commit message"

# Run security checks
pnpm audit
pnpm lint
pnpm test

# Create pull request with security checklist
```

### Security Checklist for PRs

- [ ] No hardcoded secrets or credentials
- [ ] Input validation implemented
- [ ] SQL injection prevention verified
- [ ] XSS prevention implemented
- [ ] Authentication/authorization tested
- [ ] Error messages don't leak sensitive info
- [ ] Dependencies updated and audited
- [ ] Tests include security scenarios

## Additional Resources

### Documentation
- [README.md](./README.md) - Project overview
- [CONTRIBUTING.md](./CONTRIBUTING.md) - Contribution guidelines
- [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md) - Community standards

### Security Tools
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [npm audit](https://docs.npmjs.com/cli/v8/commands/npm-audit)
- [Snyk](https://snyk.io/) - Vulnerability scanning
- [SonarQube](https://www.sonarqube.org/) - Code quality & security

### Security Training
- [OWASP WebGoat](https://owasp.org/www-project-webgoat/)
- [PortSwigger Web Security Academy](https://portswigger.net/web-security)
- [HackerOne Resources](https://www.hackerone.com/resources)

## Acknowledgments

We would like to thank the following security researchers for their contributions:

*No security researchers have been acknowledged yet.*

---

- **Last Updated**: December 23, 2025
- **Version**: 1.0.1
- **Project**: TelemetryFlow SDK

**Built with ❤️ by DevOpsCorner Indonesia**
