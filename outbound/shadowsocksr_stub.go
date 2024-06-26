//go:build !with_shadowsocksr

package outbound

import (
	"context"

	"github.com/kelleygo/sing-box/adapter"
	"github.com/kelleygo/sing-box/log"
	"github.com/kelleygo/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
)

func NewShadowsocksR(ctx context.Context, router adapter.Router, logger log.ContextLogger, tag string, options option.ShadowsocksROutboundOptions) (adapter.Outbound, error) {
	return nil, E.New("ShadowsocksR is deprecated and removed in sing-box 1.6.0")
}
