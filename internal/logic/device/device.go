package device

import "OneGfServer/internal/service"

type sDevice struct{}

func New() *sDevice {
	return &sDevice{}
}

func init() {
	service.RegisterDevice(New())
}
