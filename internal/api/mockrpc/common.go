package mockrpc

import (
	"context"
	"log"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func getHeaderMap(ctx *fiber.Ctx) map[string]string {
	allHeaders := ctx.GetReqHeaders()
	if len(allHeaders) == 0 {
		return nil
	}
	return lo.MapValues(ctx.GetReqHeaders(), func(v []string, _ string) string {
		if len(v) == 0 {
			log.Println("empty slice in http headers")
			return ""
		}
		return v[0]
	})
}

type httpRpcApi struct {
	app *fiber.App
}

func (api httpRpcApi) Listen(ctx context.Context, addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	if err := api.app.Listener(ln); err != nil {
		return err
	}
	<-ctx.Done()
	if err := ln.Close(); err != nil {
		return err
	}
	return ctx.Err()
}

type IMockApi interface {
	Register(route, method string) error
	Listen(ctx context.Context, addr string) error
}
