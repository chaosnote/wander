package game

import "os"

var (
	dc_id   = os.Getenv("DC_ID")
	dc_addr = os.Getenv("DC_ADDR")

	game_addr = os.Getenv("GAME_ADDR")

	group_id = os.Getenv("GROUP_ID")

	db_addr = os.Getenv("DB_ADDR")
	db_user = os.Getenv("DB_USER")
	db_pw   = os.Getenv("DB_PW")
	db_name = os.Getenv("DB_NAME")

	redis_addr = os.Getenv("REDIS_ADDR")
	// redis_user   = os.Getenv("REDIS_USER")
	// redis_pw     = os.Getenv("REDIS_PW")
	redis_db_idx = os.Getenv("REDIS_DB_INDEX")

	nats_addr = os.Getenv("NATS_ADDR")

	log_dir = os.Getenv("LOG_DIR")
)
