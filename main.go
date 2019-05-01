package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/aloxc/gappuser/io"
	"github.com/aloxc/gappuser/module"
	_ "github.com/aloxc/gappuser/module"
	"github.com/cyberdelia/go-metrics-graphite"
	"github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"net"
	"time"
)

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
	log.Info("password = ", request.Params["password"], user.Password)
	err := module.GetUser(&user)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("user ", user)
	user.Password = "oko"
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
	}
	return nil
}

func main() {
	srv := server.NewServer(server.WithReadTimeout(time.Duration(2)*time.Second), server.WithWriteTimeout(time.Duration(2)*time.Second))
	p := serverplugin.NewMetricsPlugin(metrics.DefaultRegistry)
	srv.Plugins.Add(p)
	srv.RegisterName("gappuser", new(GappUser), "")
	srv.Serve("tcp", ":13331")

}
func startMetrics() {
	metrics.RegisterRuntimeMemStats(metrics.DefaultRegistry)
	go metrics.CaptureRuntimeMemStats(metrics.DefaultRegistry, time.Second)

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:2003")
	go graphite.Graphite(metrics.DefaultRegistry, 1e9, "rpcx.services.host.127_0_0_1", addr)
}
