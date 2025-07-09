package queue

type sQueue struct{}

func New() *sQueue {
	return &sQueue{}
}

// func init() {

// 	service.RegisterQueue(New())
// }
