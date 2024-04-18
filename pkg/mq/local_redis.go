package mq

import (
	"context"

	"github.com/alicebob/miniredis/v2"
)

type LocalRedis struct {
	mr *miniredis.Miniredis
}

func NewLocalRedis() (*LocalRedis, error) {
	return &LocalRedis{
		mr: miniredis.NewMiniRedis(),
	}, nil
}

func (lr *LocalRedis) Listen(ctx context.Context, addr string) error {
	if err := lr.mr.StartAddr(addr); err != nil {
		return err
	}
	<-ctx.Done()
	lr.mr.Close()
	return ctx.Err()
}
