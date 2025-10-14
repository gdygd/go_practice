package api

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"order-service/internal/logger"

	"github.com/gin-gonic/gin"
)

func HttpRequest(ctx *gin.Context, payload []byte, method string, baseurl string) error {
	url := fmt.Sprintf("%s/saga/order", baseurl)

	logger.Log.Print(2, "request url : %s", url)

	// req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		logger.Log.Error("Http Request new reqeuest error.. %v", err)
		return err
	}

	// 헤더 복사 (JWT, Trace-Id, Correlation-Id 등)
	for key, values := range ctx.Request.Header {
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}

	req.Header.Set("Content-Type", "application/json")

	var client *http.Client
	client = &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableKeepAlives: true,
	}}

	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error("Http Request 호출 실패#1: %w", err)
		return fmt.Errorf("Http Request 호출 실패: %w", err)
	}
	defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	body, _ := io.ReadAll(resp.Body)
	// 	logger.Log.Error("Http Request 호출 실패#2: [%d]: %s", resp.StatusCode, string(body))
	// 	return fmt.Errorf("Http Request 호출 실패 [%d]: %s", resp.StatusCode, string(body))
	// }

	return nil
}
