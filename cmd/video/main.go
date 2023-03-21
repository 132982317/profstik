package main

import (
	"fmt"
	"github.com/132982317/profstik/cmd/video/service"
	"github.com/132982317/profstik/kitex_gen/video/videoservice"
	"github.com/132982317/profstik/pkg/utils/viper"
	"github.com/132982317/profstik/pkg/utils/zap"
	"github.com/cloudwego/kitex-examples/bizdemo/easy_note/pkg/bound"
	"github.com/cloudwego/kitex-examples/bizdemo/easy_note/pkg/middleware"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"log"
	"net"
)

var (
	config      = viper.Init("video")
	serviceName = config.Viper.GetString("server.name")
	serviceAddr = fmt.Sprintf("%s:%d", config.Viper.GetString("server.host"), config.Viper.GetInt("server.port"))
	etcdAddr    = fmt.Sprintf("%s:%d", config.Viper.GetString("etcd.host"), config.Viper.GetInt("etcd.port"))
	logger      = zap.InitLogger()
)

func main() {
	r, err := etcd.NewEtcdRegistry([]string{etcdAddr})
	if err != nil {
		log.Fatal(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", serviceAddr)
	if err != nil {
		log.Fatal(err)
	}
	svr := videoservice.NewServer(new(service.VideoServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}), // server name
		server.WithMiddleware(middleware.CommonMiddleware),                               // middleware
		server.WithMiddleware(middleware.ServerMiddleware),
		server.WithServiceAddr(addr), // address
		//server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(),                           // Multiplex
		server.WithSuite(trace.NewDefaultServerSuite()),     // tracer
		server.WithBoundHandler(bound.NewCpuLimitHandler()), // BoundHandler
		server.WithRegistry(r),                              // registry
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)
	if err := svr.Run(); err != nil {
		logger.Fatalf("%v stopped with error: %v", serviceName, err.Error())
	}
}
