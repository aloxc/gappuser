package module

import (
	"github.com/aloxc/gappuser/config"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smallnest/rpcx/log"
	"os"
	"strconv"
	"time"
)

func init() {
	dbName := os.Getenv(config.USER_MYSQL_DATABASE_NAME)
	host := os.Getenv(config.USER_MYSQL_HOST)
	user := os.Getenv(config.USER_MYSQL_USER)
	password := os.Getenv(config.USER_MYSQL_PASSWORD)
	if dbName == "" {
		dbName = config.USER_MYSQL_DATABASE_NAME_DEFAULT
	}
	if host == "" {
		host = config.USER_MYSQL_HOST_DEFAULT
	}
	if user == "" {
		user = config.USER_MYSQL_USER_DEFAULT
	}
	if password == "" {
		password = config.USER_MYSQL_PASSWORD_DEFAULT
	}

	second := os.Getenv(config.WAIT_MYSQL_SETUP_SECOND)
	var isecond int = 10
	var err error
	if second != "" {
		isecond, err = strconv.Atoi(second)
		if err != nil {
			log.Info("等待mysql初始化异常", err)
			os.Exit(0)
		}
	}

	log.Infof("等待[%d]秒后准备初始化数据库表", isecond)
	time.Sleep(time.Second * time.Duration(isecond))
	ds := user + ":" + password + "@tcp(" + host + ")/" + dbName + "?charset=utf8mb4&loc=Local"
	log.Info(ds)
	// set default database
	orm.RegisterDataBase("default", "mysql", ds, 30)
	//orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:33066)/user?charset=utf8", 30)

	orm.RegisterModel(new(User))
	// create table
	err = orm.RunSyncdb("default", true, true)
	if err != nil {
		log.Info("启动创建数据连接异常", err)
	}
}
