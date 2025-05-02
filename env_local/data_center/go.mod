module idv/chris

go 1.23.0

replace github.com/chaosnote/melody => ../../../melody

replace github.com/chaosnote/wander => ../../server

require github.com/chaosnote/wander v0.0.0-00010101000000-000000000000

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
