package module

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/aloxc/gappuser/cache"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/smallnest/rpcx/log"
	"strconv"
	"time"
)

type User_Level uint8

const (
	User_Level_OK         User_Level = iota //正常用户
	User_Level_DENY                         //被封杀用户
	User_LEVEL_BLACK_LIST                   //黑名单用户
)

type User struct {
	Id         int32
	UserName   string `orm:"unique;size(20)"`
	Password   string `orm:"size(64)"`
	Age        uint8
	Level      User_Level
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
	Version    int32
}

func GetUser(user *User) error {
	rc := cache.RedisClient.Get()
	defer rc.Close()
	reply, err := rc.Do("GET", "user:"+strconv.Itoa(int(user.Id)))
	if reply != nil {
		js, err := redis.String(reply, nil)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(js), &user)
		if err != nil {
			return err
		}
		return nil
	}
	orm.Debug = true

	orm := orm.NewOrm()
	if user.Password != "" {
		orm.QueryTable(user).Filter("id", user.Id).Filter("password", user.Password).One(user)
	} else {
		err = orm.Read(user)
	}
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(user)
	if err != nil {
		log.Info("json序列化异常", err)
		return nil
	}
	_, err = rc.Do("SET", "user:"+strconv.Itoa(int(user.Id)), string(bytes[0:]))
	if err != nil {
		log.Info("设置缓存异常", err)
		return nil
	}
	return nil
}
func InsertTestUser() {
	password := "111111"
	sum256 := sha256.Sum256([]byte(password))
	user := &User{
		Id:       1,
		UserName: "aloxc",
		Age:      12,
		Level:    0,
		Version:  1,
		Password: hex.EncodeToString(sum256[0:]),
	}
	orm := orm.NewOrm()
	orm.Insert(user)
}
func Register(user *User) error {
	orm := orm.NewOrm()
	id, err := orm.Insert(user)
	if err != nil {
		return err
	}
	user.Id = int32(id)
	log.Info("注册成功", user)
	return nil
}
func UpdateUser(user *User) error {
	orm := orm.NewOrm()
	orm.Update("")
	id, err := orm.Insert(user)
	if err != nil {
		return err
	}
	user.Id = int32(id)
	log.Info("注册成功", user)
	return nil
}
