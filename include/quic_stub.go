//go:build !with_quic

package include

import (
	"context"

	"github.com/kelleygo/sing-box/adapter"
	"github.com/kelleygo/sing-box/common/tls"
	C "github.com/kelleygo/sing-box/constant"
	"github.com/kelleygo/sing-box/option"
	"github.com/kelleygo/sing-box/transport/v2ray"
	"github.com/sagernet/sing-dns"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
)

func init() {
	dns.RegisterTransport([]string{"quic", "h3"}, func(options dns.TransportOptions) (dns.Transport, error) {
		return nil, C.ErrQUICNotIncluded
	})
	v2ray.RegisterQUICConstructor(
		func(ctx context.Context, options option.V2RayQUICOptions, tlsConfig tls.ServerConfig, handler adapter.V2RayServerTransportHandler) (adapter.V2RayServerTransport, error) {
			return nil, C.ErrQUICNotIncluded
		},
		func(ctx context.Context, dialer N.Dialer, serverAddr M.Socksaddr, options option.V2RayQUICOptions, tlsConfig tls.Config) (adapter.V2RayClientTransport, error) {
			return nil, C.ErrQUICNotIncluded
		},
	)
}
