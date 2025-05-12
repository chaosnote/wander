# 測試範例

``` sh
cd /home/chris/git/wander/env_local

sudo docker-compose up -d

sudo docker-compose down
```

``` cmd
-new_console:n:t:[DC] cmd.exe /k ""%ConEmuBaseDir%\CmdInit.cmd"&run_dc.bat"

-new_console:n:t:[Server] cmd.exe /k ""%ConEmuBaseDir%\CmdInit.cmd"&run_game.bat 0000"

-new_console:n:t:[Monkey] cmd.exe /k ""%ConEmuBaseDir%\CmdInit.cmd"&run_monkey.bat 0000"
```
