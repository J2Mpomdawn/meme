package service

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/joho/godotenv"
)

var DbEngine *xorm.Engine

func init() {
	driverName := "mysql"
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := godotenv.Load("dev.env"); err != nil {
		fmt.Println(err)
	}

	c := mysql.Config{
		DBName:    os.Getenv("DB_Name"),
		User:      os.Getenv("DB_User"),
		Passwd:    os.Getenv("DB_Passwd"),
		Addr:      os.Getenv("DB_Addr"),
		Net:       os.Getenv("DB_Net"),
		ParseTime: true,
		Collation: os.Getenv("DB_Collation"),
		Loc:       jst,
	}

	DbEngine, err = xorm.NewEngine(driverName, c.FormatDSN())
	if err != nil && err.Error() != "" {
		log.Fatal(err.Error())
	}
	DbEngine.ShowSQL(true)
	//DbEngine.SetMaxOpenConns(2)
	//DbEngine.Sync2(new(model.Value_GuildId))
	fmt.Println("init data base ok")
}
