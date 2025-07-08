package queue

import "OneGfServer/internal/service"

type sQueue struct{}

func New() *sQueue {
	return &sQueue{}
}

func init() {
	// TODO: Fix service registration - queue should not be registered as device
	// service.RegisterDevice(New())
	service.RegisterQueue(New())
}
