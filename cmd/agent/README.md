# cmd/agent

В данной директории будет содержаться код Агента, который скомпилируется в бинарное приложение.

## HOWTO
```bash
go run -ldflags "-X main.buildVersion=v1.0.1 \
    -X 'main.buildDate=$(date +%s)' \
    -X 'main.buildCommit=$(git show -s --format=%H)'" \
    cmd/agent/main.go
```
