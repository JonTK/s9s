# Enterprise Features

s9s is designed as a terminal UI interface for SLURM. For enterprise requirements, s9s leverages SLURM's native enterprise capabilities rather than reimplementing them.

## Authentication

### OAuth2/OIDC Support

s9s includes OAuth2/OIDC authentication support for integration with enterprise identity providers.

**Supported Providers:**
- Google OAuth2
- Okta (with discovery URL)
- Azure AD (with discovery URL)
- Custom OAuth2/OIDC providers

**Configuration:**

```yaml
contexts:
  - name: production
    cluster:
      endpoint: "https://slurm.example.com:6820"
      auth:
        type: oauth2
        provider: google  # or okta, azure-ad, custom
        client_id: "${OAUTH_CLIENT_ID}"
        client_secret: "${OAUTH_CLIENT_SECRET}"
        # For Okta/Azure AD/Custom:
        discovery_url: "https://your-idp.example.com/.well-known/openid-configuration"
        # Optional:
        scopes: "openid profile email"
        redirect_uri: "http://localhost:8080/callback"
```

**Setup Wizard:**

The s9s setup wizard includes OAuth2 configuration:

```bash
s9s setup
# Select option 3 for OAuth2 authentication
```

**Features:**
- OIDC discovery for automatic endpoint detection
- PKCE (Proof Key for Code Exchange) support
- Automatic token refresh
- Local callback server for authorization flow
- Support for custom scopes and redirect URIs

See [Configuration Guide](/docs/getting-started/configuration.md) for detailed authentication setup.

## Enterprise Capabilities via SLURM

For enterprise requirements beyond authentication, s9s relies on SLURM's native capabilities:

### Security & Access Control

**SLURM provides:**
- Multi-Factor Authentication (MFA) via PAM integration
- Pluggable Authentication Modules (PAM)
- Account-based access control
- Job submission policies and limits
- Resource access restrictions

**s9s integration:**
- s9s respects SLURM's authentication and authorization
- All operations are subject to SLURM's security policies
- User permissions are enforced by SLURM

### High Availability & Scalability

**SLURM provides:**
- Controller failover and redundancy
- Distributed architecture
- Multi-cluster federation
- Database redundancy (MySQL, MariaDB with replication)

**s9s integration:**
- Supports connections to highly available SLURM endpoints
- Can be configured with multiple cluster contexts
- No single point of failure when SLURM is configured for HA

### Multi-Tenancy & Resource Management

**SLURM provides:**
- Account hierarchies for organizational structure
- Fair-share scheduling across accounts/users
- Resource quotas and limits per account
- QoS (Quality of Service) policies
- Partition-based resource isolation

**s9s integration:**
- Full visibility into account hierarchies
- QoS and partition management views
- User and account resource tracking
- Reservation management

### Audit & Compliance

**SLURM provides:**
- Complete job accounting database
- Detailed audit logs
- Resource usage tracking
- Job history and provenance

**s9s integration:**
- Export job data to CSV/JSON for compliance reporting
- Real-time monitoring of resource usage
- Historical job data access

### Monitoring & Observability

**s9s provides:**
- Real-time cluster monitoring
- Job and node status visibility
- Resource utilization metrics
- Optional observability plugin for Prometheus integration

**See:** [Observability Plugin](/docs/plugins/observability.md)

## Configuration for Enterprise Environments

### Multiple Cluster Contexts

Configure multiple SLURM clusters:

```yaml
currentContext: production

contexts:
  - name: production
    cluster:
      endpoint: "https://prod-slurm.example.com:6820"
      auth:
        type: oauth2
        provider: okta
        # ... oauth config ...

  - name: development
    cluster:
      endpoint: "https://dev-slurm.example.com:6820"
      auth:
        type: slurm-token

  - name: research
    cluster:
      endpoint: "https://research-slurm.example.com:6820"
      auth:
        type: oauth2
        provider: azure-ad
        # ... oauth config ...
```

Switch between clusters:

```bash
s9s config use-context production
s9s config use-context development
```

### Secure Credential Storage

s9s supports secure credential storage through system keyrings:

```yaml
# Enable keyring storage for OAuth tokens
storage:
  backend: keyring  # or: file, memory
```

### TLS Configuration

For secure communication with SLURM REST API:

```yaml
contexts:
  - name: production
    cluster:
      endpoint: "https://slurm.example.com:6820"
      insecure: false  # Enforce TLS certificate validation
      timeout: 30s
```

## Deployment Considerations

### Container Deployment

s9s can be deployed in containerized environments:

```dockerfile
FROM alpine:latest
COPY s9s /usr/local/bin/s9s
RUN chmod +x /usr/local/bin/s9s

# Run in non-interactive mode for monitoring
ENTRYPOINT ["/usr/local/bin/s9s"]
CMD ["jobs", "--format", "json"]
```

### SSH Integration

For direct node access in enterprise environments:

```yaml
ssh:
  enabled: true
  multiplexing: true
  control_path: "/tmp/s9s-ssh-%r@%h:%p"
```

See [SSH Integration Guide](/docs/guides/ssh-integration.md) for details.

## Future Development

Additional enterprise features are under consideration. See the [specs/missing-features/](/specs/missing-features/) directory for detailed specifications of features being evaluated:

- Advanced backup and recovery capabilities
- Extended API integrations
- Enhanced multi-cluster management

For feature requests or to discuss enterprise requirements, please [open a discussion](https://github.com/jontk/s9s/discussions) or [file an issue](https://github.com/jontk/s9s/issues).

## Support

- **Community Support**: [GitHub Discussions](https://github.com/jontk/s9s/discussions)
- **Bug Reports**: [GitHub Issues](https://github.com/jontk/s9s/issues)
- **Contributing**: [Development Guide](/docs/development/contributing.md)

## Resources

- [SLURM Security Guide](https://slurm.schedmd.com/security.html)
- [SLURM High Availability](https://slurm.schedmd.com/ha.html)
- [SLURM Accounting](https://slurm.schedmd.com/accounting.html)
- [SLURM Multi-Cluster](https://slurm.schedmd.com/multi_cluster.html)
