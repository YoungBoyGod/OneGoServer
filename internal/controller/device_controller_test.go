package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YoungBoyGod/OneGoServer/internal/service"
	"github.com/gin-gonic/gin"
)

type dummyDeviceService struct{ service.DeviceService }

func TestRegisterDeviceHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	dc := NewDeviceController(&dummyDeviceService{})
	router.POST("/devices", dc.RegisterDevice)

	payload := map[string]interface{}{
		"device_id": "dev001",
		"name":      "设备A",
		"type":      "sensor",
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/devices", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		// dummy service 返回空实现，会导致内部错误，所以上下限判断
		t.Fatalf("expected 400 got %d", w.Code)
	}
}
