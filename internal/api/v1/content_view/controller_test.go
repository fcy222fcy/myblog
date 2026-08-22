package content_view

import (
	"blog/internal/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type fakeContentViewService struct {
	contentType string
	contentID   uint
	visitorKey  string
}

func (f *fakeContentViewService) Record(contentType string, contentID uint, visitorKey string) (*service.ContentViewResult, error) {
	f.contentType = contentType
	f.contentID = contentID
	f.visitorKey = visitorKey
	return &service.ContentViewResult{Counted: true, ViewCount: 12}, nil
}

func (f *fakeContentViewService) CountToday(contentType string) (int64, error) {
	return 0, nil
}

func TestControllerRecordsViewWithDerivedVisitorKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &fakeContentViewService{}
	controller := NewController(svc, "test-secret")
	router := gin.New()
	router.POST("/views", controller.Record)
	body, err := json.Marshal(map[string]any{"content_type": "article", "content_id": 7})
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/views", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Visitor-ID", "550e8400-e29b-41d4-a716-446655440000")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, "article", svc.contentType)
	require.Equal(t, uint(7), svc.contentID)
	require.Len(t, svc.visitorKey, 64)
	require.Contains(t, recorder.Body.String(), `"view_count":12`)
}

func TestControllerRejectsInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	controller := NewController(&fakeContentViewService{}, "test-secret")
	router := gin.New()
	router.POST("/views", controller.Record)
	req := httptest.NewRequest(http.MethodPost, "/views", bytes.NewBufferString(`{"content_type":"unknown","content_id":0}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}
