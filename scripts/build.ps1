# Canonical build (bun frontend + Go backend), BelowNormal priority.
$ErrorActionPreference='Stop'
(Get-Process -Id $PID).PriorityClass='BelowNormal'
bun install --frozen-lockfile
bun run build
Push-Location backend-go
$env:CGO_ENABLED='0'
go vet ./...
go build -o jules-backend.exe .
Pop-Location
