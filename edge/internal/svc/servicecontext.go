package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"zeroim/edge/internal/config"
	"zeroim/imrpc/imrpc_client"
)

type ServiceContext struct {
	Config config.Config
	IMRpc  imrpc_client.Imrpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	client := zrpc.MustNewClient(c.IMRpc)
	return &ServiceContext{
		Config: c,
		IMRpc:  imrpc_client.NewImrpc(client),
	}
}
