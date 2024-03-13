package service

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"golang.org/x/crypto/ssh"

	"meme/model"
)

// DataBase Instance
var DbEngine *xorm.Engine

// SSH Instance
var SshClient *ssh.Client

// connect SSH and DB
func Connect() {
	sshConf := &model.SSH{
		Key:  os.Getenv("SSH_Key"),
		Host: os.Getenv("SSH_Host"),
		Port: os.Getenv("SSH_Port"),
		User: os.Getenv("SSH_User"),
	}

	dbConf := &model.DB{
		Host:      os.Getenv("DB_Host"),
		Port:      os.Getenv("DB_Port"),
		User:      os.Getenv("DB_User"),
		Password:  os.Getenv("DB_Password"),
		DBName:    os.Getenv("DB_Name"),
		Collation: os.Getenv("DB_Collation"),
	}

	//get SSH Instance
	SshClient, err := newSSH(sshConf)

	if err != nil {
		LogPrintln("red", "newSSH", err)

		//permission ssh
		SSHPermission()

		//You will be able to connect in about 5 minutes after registration.
		time.Sleep(6 * time.Minute)

		//connection retry
		SshClient, err = newSSH(sshConf)

		if err != nil {
			return
		}
	}

	FmtPrintln("blue", "init ssh ok")

	//defer SshClient.Close()

	//get DB Instance
	DbEngine, err = newDB(dbConf, SshClient)

	if err != nil && err.Error() != "" {
		LogPrintln("red", "newDB", err)
	}

	FmtPrintln("blue", "init data base ok")
}

// disconnect SSH and DB
func DisConnect() {
	//close DB
	DbEngine.Close()

	//close SSH
	SshClient.Close()
}

// create DB Instance
func newDB(conf *model.DB, sshc *ssh.Client) (*xorm.Engine, error) {
	driverName := "mysql"
	mysqlNet := "tcp"

	if sshc != nil {
		mysqlNet = "mysql+tcp"

		mysql.RegisterDialContext(mysqlNet, func(ctx context.Context, addr string) (net.Conn, error) {
			return sshc.Dial("tcp", addr)
		})
	}

	jst, err := time.LoadLocation("Asia/Tokyo")

	if err != nil {
		return nil, err
	}

	c := mysql.Config{
		DBName:               conf.DBName,
		User:                 conf.User,
		Passwd:               conf.Password,
		Addr:                 net.JoinHostPort(conf.Host, conf.Port),
		Net:                  mysqlNet,
		ParseTime:            true,
		Collation:            conf.Collation,
		Loc:                  jst,
		AllowNativePasswords: true,
	}

	//connection process
	DbEngine, err := xorm.NewEngine(driverName, c.FormatDSN())

	if err != nil {
		LogPrint("red", "NewEngine", err)
	}

	//DbEngine.ShowSQL(true)
	DbEngine.SetMaxOpenConns(2)

	return DbEngine, err
}

// create SSH Instance
func newSSH(conf *model.SSH) (*ssh.Client, error) {
	sshConf := &ssh.ClientConfig{
		User: conf.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(conf.Key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	//connection process
	return ssh.Dial("tcp", net.JoinHostPort(conf.Host, conf.Port), sshConf)
}
