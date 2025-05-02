package datacenter

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/chaosnote/wander/model/member"
	"github.com/chaosnote/wander/utils"
)

type DCStore interface {
	Close()
}

//-----------------------------------------------

type dc_store struct {
	utils.LogStore

	db_store     *sql.DB
	player_store map[string]*member.Player
}

//-----------------------------------------------

func NewDCStore(logger utils.LogStore) DCStore {
	ds := &dc_store{
		LogStore:     logger,
		player_store: make(map[string]*member.Player),
	}

	var e error

	// 例 : "user:password@tcp(ip)?parseTime=true/dbname"
	cmd := fmt.Sprintf(`%s:%s@tcp(%s)/%s?parseTime=true`, db_user, db_pw, db_addr, db_name)
	ds.Debug(utils.LogFields{"cmd": cmd})
	ds.db_store, e = sql.Open("mysql", cmd)
	if e != nil {
		panic(e)
	}
	e = ds.db_store.Ping()
	if e != nil {
		panic(e)
	}
	ds.db_store.SetMaxOpenConns(100)          // Limit to N open connections
	ds.db_store.SetMaxIdleConns(10)           // Keep up to N idle connections
	ds.db_store.SetConnMaxLifetime(time.Hour) // Reuse connections for at most N

	utils.RSAInit("./asset/_rsa/pem.key", 512, true)

	ds.initRouter()

	return ds
}

//-----------------------------------------------

func (ds *dc_store) Close() {
	ds.db_store.Close()
}
