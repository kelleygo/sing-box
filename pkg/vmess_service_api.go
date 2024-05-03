package pkg

import "context"

type Tun2VmessService interface {
	Start(ctx context.Context) error
	Create(ctx context.Context) error
	Close() error
}
