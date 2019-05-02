package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/aloxc/gappuser/config"
	"github.com/aloxc/gappuser/io"
	"github.com/aloxc/gappuser/module"
	_ "github.com/aloxc/gappuser/module"
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"net"
	"os"
	"strconv"
	"time"
)

func init() {
	module.InsertTestUser()
}

type GappUser struct {
}

func (this *GappUser) getUser(ctx context.Context, request *io.Request, response *io.Response) error {
	var user module.User
	user = module.User{
		Id: request.Params["userId"].(int32),
	}
	if request.Params["password"] != nil {
		sum256 := sha256.Sum256([]byte(request.Params["password"].(string)))
		user.Password = hex.EncodeToString(sum256[0:])
	}
	err := module.GetUser(&user)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("user ", user)
	name, err := os.Hostname()
	if err != nil {
		log.Info("读取本机名异常", err)
		name = err.Error()
	}
	name = "...host = [" + name + "]"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Info("读取本机ip异常", err)
		name += "  " + err.Error()
	}
	name = name + "ip = ["
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				name += "(" + ipnet.IP.String() + ")"
			}
		}
	}
	name += "]"
	user.Password = "oko" + name
	response.Code = 0
	response.Message = "正常请求"
	response.Data = user
	return nil
}
func (this *GappUser) register(ctx context.Context, request *io.Request, response *io.Response) error {
	var user module.User
	user = module.User{
		UserName: request.Params["userName"].(string),
		Password: request.Params["password"].(string),
		Level:    module.User_Level(request.Params["level"].(int64)),
		Version:  1,
	}
	sum256 := sha256.Sum256([]byte(request.Params["password"].(string)))
	user.Password = hex.EncodeToString(sum256[0:])

	err := module.Register(&user)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("user ", user)
	response.Code = 0
	response.Message = "正常请求"
	response.Data = user
	return nil
}
func (this *GappUser) updateUser(ctx context.Context, request *io.Request, response *io.Response) error {
	var user module.User
	user = module.User{
		UserName: request.Params["userName"].(string),
		Password: request.Params["password"].(string),
		Level:    module.User_Level(request.Params["level"].(int64)),
		Version:  1,
	}
	sum256 := sha256.Sum256([]byte(request.Params["password"].(string)))
	user.Password = hex.EncodeToString(sum256[0:])

	err := module.Register(&user)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("user ", user)
	response.Code = 0
	response.Message = "正常请求"
	response.Data = user
	return nil
}

func (this *GappUser) Execute(ctx context.Context, request *io.Request, response *io.Response) error {
	bytes, _ := json.Marshal(request)
	log.Info("请求", string(bytes[0:]))
	switch request.Method {
	case "getUser":
		return this.getUser(ctx, request, response)
	case "register":
		return this.register(ctx, request, response)
	case "updateUser":
		return this.updateUser(ctx, request, response)
	}
	return nil
}

func main() {
	port := os.Getenv(config.SERVER_PORT)
	if port == "" {
		port = strconv.Itoa(config.SERVER_PORT_DEFAULT)
	}
	srv := server.NewServer(server.WithReadTimeout(time.Duration(2)*time.Second), server.WithWriteTimeout(time.Duration(2)*time.Second))
	p := serverplugin.NewMetricsPlugin(metrics.DefaultRegistry)
	srv.Plugins.Add(p)
	srv.RegisterName("gappuser", new(GappUser), "")
	err := srv.Serve("tcp", ":"+port)
	if err != nil {
		log.Error("rpcx服务无法启动", err)
		os.Exit(1)
	}
	log.Info("服务已启动，监听端口[", port, "]")
}
