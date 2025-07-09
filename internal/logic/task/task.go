package task

type sTask struct{}

func New() *sTask {
	return &sTask{}
}

// func init() {
// 	service.RegisterTask(New())
// }
