@echo off

@REM 遊戲共用-例: .\build-game.bat "0000"
@REM 切換到批次檔所在的目錄

if "%1"=="" (
    echo 錯誤：未設置 game_id
    exit /b 1
)

cd /d %~dp0

set out_dir="../env_local/game_%1/protobuf"

if exist %out_dir% (
    echo remove dir dist
    rmdir /q /s %out_dir%
) 

mkdir %out_dir%

set PATH=%PATH%;./bin/

protoc --proto_path="./src" --go_out=%out_dir% ^
game_%1.proto