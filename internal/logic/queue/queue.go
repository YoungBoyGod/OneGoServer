package queue

import (
	"golang.org/x/net/context"
)

type sQueue struct{}

func New() *sQueue {
	return &sQueue{}
}

func init() {
	service.RegisterQueue(New())
}

func (s *sQueue) GetQueueInstance(ctx context.Context) (*queue.Client, error) {
}
