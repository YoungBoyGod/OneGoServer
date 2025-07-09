package device

type sDevice struct{}

func New() *sDevice {
	return &sDevice{}
}

// func init() {
// 	service.RegisterDevice(New())
// }
