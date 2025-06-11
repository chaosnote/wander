@echo off

set vb_addr=192.168.0.236

@REM DC 伺服器
set DC_ID=28c2f2c2-590f-442d-bf80-7e6bbeb4c471
set DC_ADDR=:8080

@REM GAME
set GAME_ADDR=:8081
set GROUP_ID=00

@REM DB
set DB_ADDR=%vb_addr%:3306
set DB_USER=chris
set DB_PW=123456
set DB_NAME=game_dev

@REM MONGO
set MONGO_ADDR=%vb_addr%:27017
set MONGO_USER=admin
set MONGO_PW=password

@REM NATS
set NATS_ADDR=%vb_addr%:4222

@REM REDIS
set REDIS_ADDR=%vb_addr%:6379
set REDIS_DB_INDEX=0

@REM LOG_DIR
set LOG_DIR=./logs/

if "%1"=="" (
    echo 錯誤：未設置 game_id
    exit /b 1
)

set proj=%cd%
cd %proj%/game_%1
go mod tidy

echo cls^&go run -race ./server/. -log_mode 2