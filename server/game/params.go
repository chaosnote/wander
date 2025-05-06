package game

import "flag"

var (
	GAME_ID  = flag.String("game_id", "0000", "game or map id")
	LOG_MODE = flag.Int("log_mode", 1, "{0:Console,1:File,2:max,3:elk[not yet]...}")
)
