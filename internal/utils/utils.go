package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Envelope: map[string]interface{} — interface{} 是空接口，表示任意类型（等价于 Go 1.18+ 的 any）。
// 用于构造灵活的 JSON 响应，如 Envelope{"data": workout, "count": 42}。
type Envelope map[string]interface{}

func WriteJsonResponse(w http.ResponseWriter, statusCode int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	// 重要顺序：必须先 Header().Set()，再 WriteHeader()。因为 WriteHeader() 一旦调用，响应头就被发送出去了，之后再 Set 就无效了。所以你的代码中第 18-19 行的顺序是正确的。
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(js)
	if err != nil {
		return err
	}
	return nil
}

func ReadIdParam(r *http.Request, paramName string) (int64, error) {
	paramsID := chi.URLParam(r, paramName)
	if paramsID == "" {
		return 0, fmt.Errorf("id parameter is required")
	}
	id, err := strconv.ParseInt(paramsID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id parameter")
	}
	return id, nil
}
