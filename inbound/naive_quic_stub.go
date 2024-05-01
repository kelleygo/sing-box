//go:build !with_quic

package inbound

import (
	C "github.com/kelleygo/sing-box/constant"
)

func (n *Naive) configureHTTP3Listener() error {
	return C.ErrQUICNotIncluded
}
