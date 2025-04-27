# identity

## Install

1) Code generation
```bash
go generate
```

2) Run the app
```bash
set -a && source .env && set +a  # load env vars
go run cmd/identity/main.go
```