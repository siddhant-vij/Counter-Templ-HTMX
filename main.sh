rm templates/*_templ.go

templ generate
go mod tidy
go run cmd/main.go