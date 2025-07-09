// package queue

// import (
// 	"context"
// 	"sort"

// 	"github.com/gogf/gf/v2/os/gtime"

// 	queue "OneGfServer/internal/model/queue"
// )

// // ===============================
// // 队列排序算法业务逻辑
// // ===============================

// // SortQueueTasks 队列任务排序
// func (s *sQueue) SortQueueTasks(ctx context.Context, input *queue.SortQueueTasksInput) (*queue.SortQueueTasksOutput, error) {
// 	var sortedTasks []map[string]interface{}

// 	switch input.QueueType {
// 	case "fifo":
// 		fifoInput := &queue.SortFIFOInput{Tasks: input.Tasks}
// 		fifoOutput := s.sortFIFO(fifoInput)
// 		sortedTasks = fifoOutput.SortedTasks
// 	case "lifo":
// 		lifoInput := &queue.SortLIFOInput{Tasks: input.Tasks}
// 		lifoOutput := s.sortLIFO(lifoInput)
// 		sortedTasks = lifoOutput.SortedTasks
// 	case "priority":
// 		priorityInput := &queue.SortByPriorityInput{Tasks: input.Tasks}
// 		priorityOutput := s.sortByPriority(priorityInput)
// 		sortedTasks = priorityOutput.SortedTasks
// 	case "round_robin":
// 		roundRobinInput := &queue.SortRoundRobinInput{Tasks: input.Tasks}
// 		roundRobinOutput := s.sortRoundRobin(roundRobinInput)
// 		sortedTasks = roundRobinOutput.SortedTasks
// 	case "weighted":
// 		weightedInput := &queue.SortByWeightInput{Tasks: input.Tasks}
// 		weightedOutput := s.sortByWeight(weightedInput)
// 		sortedTasks = weightedOutput.SortedTasks
// 	default:
// 		// 默认FIFO
// 		fifoInput := &queue.SortFIFOInput{Tasks: input.Tasks}
// 		fifoOutput := s.sortFIFO(fifoInput)
// 		sortedTasks = fifoOutput.SortedTasks
// 	}

// 	return &queue.SortQueueTasksOutput{
// 		SortedTasks: sortedTasks,
// 	}, nil
// }

// // sortFIFO 先进先出排序
// func (s *sQueue) sortFIFO(input *queue.SortFIFOInput) *queue.SortFIFOOutput {
// 	sort.Slice(input.Tasks, func(i, j int) bool {
// 		timeI, _ := input.Tasks[i]["created_at"].(*gtime.Time)
// 		timeJ, _ := input.Tasks[j]["created_at"].(*gtime.Time)

// 		if timeI == nil || timeJ == nil {
// 			return false
// 		}

// 		return timeI.Before(timeJ)
// 	})
// 	return &queue.SortFIFOOutput{
// 		SortedTasks: input.Tasks,
// 	}
// }

// // sortLIFO 后进先出排序
// func (s *sQueue) sortLIFO(input *queue.SortLIFOInput) *queue.SortLIFOOutput {
// 	sort.Slice(input.Tasks, func(i, j int) bool {
// 		timeI, _ := input.Tasks[i]["created_at"].(*gtime.Time)
// 		timeJ, _ := input.Tasks[j]["created_at"].(*gtime.Time)

// 		if timeI == nil || timeJ == nil {
// 			return false
// 		}

// 		return timeI.After(timeJ)
// 	})
// 	return &queue.SortLIFOOutput{
// 		SortedTasks: input.Tasks,
// 	}
// }

// // sortByPriority 按优先级排序
// func (s *sQueue) sortByPriority(input *queue.SortByPriorityInput) *queue.SortByPriorityOutput {
// 	sort.Slice(input.Tasks, func(i, j int) bool {
// 		priorityI, _ := input.Tasks[i]["priority"].(int)
// 		priorityJ, _ := input.Tasks[j]["priority"].(int)

// 		// 优先级高的在前
// 		if priorityI != priorityJ {
// 			return priorityI > priorityJ
// 		}

// 		// 优先级相同时，按创建时间排序
// 		timeI, _ := input.Tasks[i]["created_at"].(*gtime.Time)
// 		timeJ, _ := input.Tasks[j]["created_at"].(*gtime.Time)

// 		if timeI == nil || timeJ == nil {
// 			return false
// 		}

// 		return timeI.Before(timeJ)
// 	})
// 	return &queue.SortByPriorityOutput{
// 		SortedTasks: input.Tasks,
// 	}
// }

// // sortRoundRobin 轮询排序
// func (s *sQueue) sortRoundRobin(input *queue.SortRoundRobinInput) *queue.SortRoundRobinOutput {
// 	// 按任务类型分组
// 	typeGroups := make(map[string][]map[string]interface{})
// 	for _, task := range input.Tasks {
// 		taskType, _ := task["type"].(string)
// 		typeGroups[taskType] = append(typeGroups[taskType], task)
// 	}

// 	// 轮询排序
// 	var result []map[string]interface{}
// 	indices := make(map[string]int)

// 	for len(result) < len(input.Tasks) {
// 		for taskType, typeTasks := range typeGroups {
// 			if indices[taskType] < len(typeTasks) {
// 				result = append(result, typeTasks[indices[taskType]])
// 				indices[taskType]++
// 			}
// 		}
// 	}

// 	return &queue.SortRoundRobinOutput{
// 		SortedTasks: result,
// 	}
// }

// // sortByWeight 按权重排序
// func (s *sQueue) sortByWeight(input *queue.SortByWeightInput) *queue.SortByWeightOutput {
// 	sort.Slice(input.Tasks, func(i, j int) bool {
// 		weightInputI := &queue.CalculateTaskWeightInput{Task: input.Tasks[i]}
// 		weightInputJ := &queue.CalculateTaskWeightInput{Task: input.Tasks[j]}

// 		weightI := s.calculateTaskWeight(weightInputI)
// 		weightJ := s.calculateTaskWeight(weightInputJ)

// 		return weightI.Weight > weightJ.Weight
// 	})
// 	return &queue.SortByWeightOutput{
// 		SortedTasks: input.Tasks,
// 	}
// }

// // calculateTaskWeight 计算任务权重
// func (s *sQueue) calculateTaskWeight(input *queue.CalculateTaskWeightInput) *queue.CalculateTaskWeightOutput {
// 	weight := 1.0

// 	// 优先级权重 (权重: 40%)
// 	if priority, ok := input.Task["priority"].(int); ok {
// 		weight += float64(priority) * 0.4
// 	}

// 	// 任务类型权重 (权重: 30%)
// 	if taskType, ok := input.Task["type"].(string); ok {
// 		switch taskType {
// 		case "urgent":
// 			weight += 3.0
// 		case "high":
// 			weight += 2.0
// 		case "normal":
// 			weight += 1.0
// 		case "low":
// 			weight += 0.5
// 		}
// 	}

// 	// 等待时间权重 (权重: 20%)
// 	if createdAt, ok := input.Task["created_at"].(*gtime.Time); ok {
// 		waitTime := gtime.Now().Sub(createdAt)
// 		weight += waitTime.Hours() * 0.2
// 	}

// 	// 资源需求权重 (权重: 10%)
// 	if resourceDemand, ok := input.Task["resource_demand"].(float64); ok {
// 		weight += resourceDemand * 0.1
// 	}

// 	return &queue.CalculateTaskWeightOutput{
// 		Weight: weight,
// 	}
// }
