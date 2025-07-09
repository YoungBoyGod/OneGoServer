package system

type sSystem struct{}

func New() *sSystem {
	return &sSystem{}
}

// func init() {
// 	service.RegisterSystem(New())
// }
