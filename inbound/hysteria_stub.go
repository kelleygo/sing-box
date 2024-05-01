//go:build !with_quic

package inbound

import (
	"context"

	"github.com/kelleygo/sing-box/adapter"
	C "github.com/kelleygo/sing-box/constant"
	"github.com/kelleygo/sing-box/log"
	"github.com/kelleygo/sing-box/option"
)

func NewHysteria(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.HysteriaInboundOptions) (adapter.Inbound, error) {
	return nil, C.ErrQUICNotIncluded
}

func NewHysteria2(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.Hysteria2InboundOptions) (adapter.Inbound, error) {
	return nil, C.ErrQUICNotIncluded
}
