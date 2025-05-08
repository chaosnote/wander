@echo off

@REM 本地伺服器
set DC_ID=28c2f2c2-590f-442d-bf80-7e6bbeb4c471
set DC_ADDR=:8080
set USE_GUEST=1

@REM DB
set DB_ADDR=192.168.0.236:3306
set DB_USER=chris
set DB_PW=123456
set DB_NAME=game_dev

@REM NATS
set NATS_ADDR=192.168.0.236:4222

@REM REDIS
set REDIS_ADDR=192.168.0.236:6379
set REDIS_DB_INDEX=0

go mod tidy
