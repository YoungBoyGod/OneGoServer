package user

type sUser struct{}

func New() *sUser {
	return &sUser{}
}

// func init() {
// 	service.RegisterUser(New())
// }
