package admin

import "flag"

var (
	LOG_MODE = flag.Int("log_mode", 1, "{0:Console,1:File,2:mix,3:elk[not yet]...}")
)
