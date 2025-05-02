module idv/chris

go 1.23.0

toolchain go1.23.8

replace github.com/chaosnote/melody => ../../../melody

replace github.com/chaosnote/wander => ../../server

require (
	github.com/chaosnote/melody v0.0.0-00010101000000-000000000000
	github.com/chaosnote/wander v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/go-resty/resty/v2 v2.16.5 // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/nxadm/tail v1.4.11 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
