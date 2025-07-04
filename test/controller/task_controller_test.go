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

type dummyTaskService struct{}

func (d *dummyTaskService) CreateTask(_, _ interface{}) (*interface{}, error) { return nil, nil }

func TestCreateTaskHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	tc := controller.NewTaskController(&dummyTaskService{})
	r.POST("/tasks", tc.CreateTask)
	body, _ := json.Marshal(map[string]interface{}{"name": "t1", "type": "backup"})
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("unexpected status %d", w.Code)
	}
}
