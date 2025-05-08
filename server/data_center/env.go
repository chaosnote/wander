package datacenter

import "os"

var (
	dc_id   = os.Getenv("DC_ID")
	dc_addr = os.Getenv("DC_ADDR")

	use_guest = os.Getenv("USE_GUEST")

	db_addr = os.Getenv("DB_ADDR")
	db_user = os.Getenv("DB_USER")
	db_pw   = os.Getenv("DB_PW")
	db_name = os.Getenv("DB_NAME")

	nats_addr = os.Getenv("NATS_ADDR")

	redis_addr   = os.Getenv("REDIS_ADDR")
	redis_db_idx = os.Getenv("REDIS_DB_INDEX")
)
