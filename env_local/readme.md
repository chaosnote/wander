# 測試範例

``` sh
cd /home/chris/git/wander/env_local

sudo docker-compose up -d

sudo docker-compose down
```

``` cmd
-new_console:n:t:[DC] cmd.exe /k ""%ConEmuBaseDir%\CmdInit.cmd"&cd ./data_center&env.bat&echo cls^&go run ."

-new_console:n:t:[Game_0000] cmd.exe /k ""%ConEmuBaseDir%\CmdInit.cmd"&cd ./game_0000&env.bat&echo cls^&go run -race . -log_mode 2"

-new_console:n:t:[Monkey] cmd.exe /k ""%ConEmuBaseDir%\CmdInit.cmd"&cd ./monkey&echo cls^&go run -race ."
```
