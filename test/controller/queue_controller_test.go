// 单元测试：QueueController EnqueueTask
// 使用 gin.TestMode 与 dummyQueueService（满足 QueueService 接口）
// 验证路由处理正常返回 200

package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YoungBoyGod/OneGoServer/internal/controller"
	"github.com/gin-gonic/gin"
)

// dummyQueueService 满足 QueueService 接口，直接返回 nil
type dummyQueueService struct{}

func TestEnqueueTaskHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	qc := controller.NewQueueController(&dummyQueueService{})
	r.POST("/device-queues/enqueue", qc.EnqueueTask)

	payload, _ := json.Marshal(map[string]interface{}{"device_id": 1, "task_id": 2, "priority": 5})
	req, _ := http.NewRequest(http.MethodPost, "/device-queues/enqueue", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", w.Code)
	}
}
