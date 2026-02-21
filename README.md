# deploy CLI

Cross-platform deployment CLI that targets remote Linux hosts over SSH using a clean, layered architecture.

## Features
- Commands: `init`, `push`, `rollback`, `status`, `logs`
- Supports Docker compose, Node/pm2, Laravel/PHP, Python, and static sites
- Remote backups and rollback with lock protection
- Stream deployment logs from the VPS

## Project layout (Clean Architecture)
```
cmd/deploy/main.go          // CLI entrypoint
internal/domain             // Enterprise models and errors
internal/application        // Use cases (init, push, rollback, status, logs)
internal/infrastructure     // SSH, remote executor, detectors, logging
internal/interfaces/cli     // Command parsing and config loader
configs/config.yaml         // Sample configuration
```

## Configuration
Set `DEPLOY_CONFIG` or use the default `configs/config.yaml`:
```yaml
ssh:
  host: example.com
  user: deploy
  port: 22
  privateKeyPath: ~/.ssh/id_rsa
  password: ""
  knownHostsPath: ~/.ssh/known_hosts
```

`knownHostsPath` must point to a valid `known_hosts` file; the CLI refuses to connect without host key verification.

## Usage
```bash
# initialize remote paths and bare repo
deploy init myapp

# deploy latest main branch
deploy push myapp

# rollback to latest or specific backup
deploy rollback myapp
deploy rollback -backup myapp-17170000.tgz myapp

# check remote status
deploy status myapp

# stream logs (tail -f)
deploy logs -f -n 200 myapp
```

Build locally with Go:
```bash
go build ./cmd/deploy
```
