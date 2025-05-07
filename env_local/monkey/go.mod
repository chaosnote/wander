module idv/chris

go 1.23.0

replace github.com/chaosnote/melody => ../../../melody

replace github.com/chaosnote/wander => ../../server

require (
	github.com/chaosnote/melody v0.0.0-00010101000000-000000000000
	github.com/chaosnote/wander v0.0.0-00010101000000-000000000000
	github.com/go-resty/resty/v2 v2.16.5
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
