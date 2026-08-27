package interceptors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/thimira/production-tracer/internal/exceptions"

	"github.com/gin-gonic/gin"
	"github.com/rashintha/logger"
)

const (
	Origin  byte = 0x0
	Replace byte = 0x1
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status byte
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	if r.status == Replace {
		r.body.Write(b)
		return r.ResponseWriter.Write(b)
	}
	return r.body.Write(b)
}

type OrderedResponse struct {
	Response  int         `json:"response"`
	Path      string      `json:"path"`
	Timestamp string      `json:"timestamp"`
	Error     *int        `json:"error,omitempty"`
	Message   interface{} `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func Interceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Next()
			return
		}
		wb := &responseBodyWriter{
			body:           &bytes.Buffer{},
			ResponseWriter: c.Writer,
			status:         Origin,
		}
		c.Writer = wb
		c.Next()
		contentType := c.Writer.Header().Get("Content-Type")
		if shouldSkipInterceptor(contentType) {
			wb.status = Replace
			wb.ResponseWriter.Write(wb.body.Bytes())
			return
		}
		wb.status = Replace
		var (
			responseData interface{}
			errorCode    *int
			errorMessage interface{}
		)
		statusCode := wb.Status()
		originBytes := wb.body.Bytes()
		if strings.HasPrefix(contentType, "application/json") {
			var jsonData interface{}
			if err := json.Unmarshal(originBytes, &jsonData); err != nil {
				code := 400
				errorCode = &code
				errorMessage = "Invalid JSON response"
			} else {
				if mapData, ok := jsonData.(map[string]interface{}); ok {
					if errMsg, hasError := mapData["error_message"]; hasError {
						code := 400
						errorCode = &code
						errorMessage = errMsg
						statusCode = 400
					} else {
						if dataVal, hasData := mapData["data"]; hasData {
							if _, hasPagination := mapData["pagination"]; hasPagination {
								responseData = mapData
							} else {
								responseData = dataVal
							}
						} else {
							responseData = mapData
						}
					}
				} else {
					responseData = jsonData
				}
			}

		} else if strings.HasPrefix(contentType, "text/plain") {
			responseData = string(originBytes)
		} else {
			code := 404
			errorCode = &code
			errorMessage = exceptions.ResourceNotFound.Error()
			statusCode = 404
		}
		if responseData == nil && errorCode == nil {
			code := 404
			errorCode = &code
			errorMessage = "Records not found"
			statusCode = 404
		}
		responseData = removeProtectedFields(responseData, []string{"password"})
		result := OrderedResponse{
			Response:  statusCode,
			Path:      c.Request.URL.Path,
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			Error:     errorCode,
			Message:   errorMessage,
			Data:      responseData,
		}
		wb.body = &bytes.Buffer{}
		c.JSON(statusCode, result)
	}
}

func shouldSkipInterceptor(contentType string) bool {
	skipTypes := []string{
		"application/pdf",
		"application/octet-stream",
		"image/",
		"video/",
		"audio/",
		"application/zip",
		"application/x-zip-compressed",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument",
	}
	for _, t := range skipTypes {
		if strings.HasPrefix(contentType, t) {
			return true
		}
	}
	return false
}

func removeProtectedFields(data interface{}, fields []string) interface{} {
	if data == nil {
		return nil
	}

	b, err := json.Marshal(data)
	if err != nil {
		return data
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return data
	}
	for _, f := range fields {
		delete(m, f)
	}
	return m
}

func Log() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Defaultln(fmt.Sprintf("ip: %s - time: %s | method: %s | path: %s | proto: %s | status: %d | latency: %s | userAgent: %s | error: %s",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage))

		return ""
	})
}
