@echo off

@REM DC 伺服器
set DC_ID=28c2f2c2-590f-442d-bf80-7e6bbeb4c471
set DC_ADDR=:8080

@REM GAME
set GAME_ADDR=:8081
set GROUP_ID=00

@REM DB
set DB_ADDR=192.168.0.236:3306
set DB_USER=chris
set DB_PW=123456
set DB_NAME=game_dev

@REM REDIS
set REDIS_ADDR=192.168.0.236:6379
@REM set REDIS_USER=chris
@REM set REDIS_PW=123456
set REDIS_DB_INDEX=0

@REM LOG_DIR
set LOG_DIR=./logs/

go mod tidy
