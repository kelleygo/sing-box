//go:build !with_quic

package outbound

import (
	"context"

	"github.com/kelleygo/sing-box/adapter"
	C "github.com/kelleygo/sing-box/constant"
	"github.com/kelleygo/sing-box/log"
	"github.com/kelleygo/sing-box/option"
)

func NewHysteria(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.HysteriaOutboundOptions) (adapter.Outbound, error) {
	return nil, C.ErrQUICNotIncluded
}

func NewHysteria2(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.Hysteria2OutboundOptions) (adapter.Outbound, error) {
	return nil, C.ErrQUICNotIncluded
}
