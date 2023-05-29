# Snippetbox
Pastebin clone using GO.

## Create TLS Certificate
```
cd tls
go run ~YOUR GO DIRECTORY~/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

## Start the application
```
go run ./cmd/web/
```
