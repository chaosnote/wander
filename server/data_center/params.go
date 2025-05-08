package datacenter

import "flag"

var (
	LOG_MODE = flag.Int("log_mode", 1, "{0:Console,1:File,2:max,3:elk[not yet]...}")
)
