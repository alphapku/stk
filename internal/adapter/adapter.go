package adapter

import "context"

type Adapter interface {
	Start(ctx context.Context, dataChan chan interface{} /*cfg*/) (<-chan struct{}, error)
	Close(ctx context.Context)
}
