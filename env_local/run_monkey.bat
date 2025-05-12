@echo off

if "%1"=="" (
    echo 錯誤：未設置 game_id
    exit /b 1
)

set proj=%cd%
cd %proj%/game_%1

echo cls^&go run -race ./monkey/.