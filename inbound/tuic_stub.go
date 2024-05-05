//go:build !with_quic

package inbound

import (
	"context"

	"github.com/kelleygo/sing-box/adapter"
	C "github.com/kelleygo/sing-box/constant"
	"github.com/kelleygo/sing-box/log"
	"github.com/kelleygo/sing-box/option"
)

func NewTUIC(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.TUICInboundOptions) (adapter.Inbound, error) {
	return nil, C.ErrQUICNotIncluded
}
