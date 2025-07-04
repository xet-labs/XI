package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"xi/app/lib"
	"xi/conf"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBService struct {
	once sync.Once
}

var DB = &DBService{}

func init() {
	DB.Init()
	lib.DB.RegisterLazyInit(DB.Init)
	lib.Redis.RegisterLazyInit(DB.Init)
}

func (d *DBService) preInit() {
	// Set global Redis and DB defaults
	lib.DB.SetDefault(conf.DbConf.DefDb)
	lib.Redis.SetCtx(context.Background())
	lib.Redis.SetDefault(conf.DbConf.RedisDefRdb)
	lib.Redis.SetPrefix(conf.DbConf.RedisPrefix)
}

func (d *DBService) postInit() {
}

// Initializes all DBs and Redis clients (forced)
func (d *DBService) InitForce() {

	d.preInit()

	for profile, c := range conf.DB {
		if !c.Enable {
			continue
		}

		// Fallback for DB credentials
		if c.User == "" {
			c.User = c.Database + "_u"
		}
		if c.Pass == "" {
			c.Pass = lib.Env.Get("DB_PASS")
		}

		switch c.Driver {
		case "mysql", "mariadb":
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
				c.User, c.Pass, c.Host, c.Port, c.Database, c.Charset)
			dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatalf("❌ Could not connect to DB '%s': %v", profile, err)
			}
			lib.DB.SetCli(profile, dbConn)
			log.Printf("✅ DB connected '%s' (MySQL)", profile)

		case "sqlite":
			dbConn, err := gorm.Open(sqlite.Open(c.Database), &gorm.Config{})
			if err != nil {
				log.Fatalf("❌ Could not connect to DB '%s': %v", profile, err)
			}
			lib.DB.SetCli(profile, dbConn)
			log.Printf("✅ DB connected '%s' (SQLite)", profile)

		case "redis":
			rdb := redis.NewClient(&redis.Options{
				Addr:     c.Host + ":" + c.Port,
				Password: c.Pass,
				DB:       c.RedisDB,
			})
			if err := rdb.Ping(context.Background()).Err(); err != nil {
				log.Fatalf("❌ Could not connect to Redis '%s': %v", profile, err)
			}
			lib.Redis.SetCli(profile, rdb)
			log.Printf("✅ DB connected '%s' (Redis)", profile)

		default:
			log.Printf("⚠️  Unsupported DB driver '%s' for DB '%s'", c.Driver, profile)
		}
	}
	
	d.postInit()
}

// Init initializes DBs once
func (d *DBService) Init() {
	d.once.Do(d.InitForce)
}
