package api

import (
	"encoding/json"
	"fmt"
	"order-delivery-saga/internal/logger"

	"github.com/gdygd/goglib"
	"github.com/gin-gonic/gin"
)

const BASE_URL = "http://10.1.0.119:9080"

func callOrderServiceConfirm(orderId int) error {

	return nil
}

func callOrderServiceCancel(orderId int) error {

	return nil
}

func callDeliveryService(ctx *gin.Context, orderId int) (int, error) {
	var delivery deliveryRequest
	delivery.OrderId = orderId
	delivery.Address = "gunpo1"
	url := fmt.Sprintf("%s/delivery/request", BASE_URL)
	// statuscode := 0

	logger.Log.Print(2, "call delivery url : %s", url)

	// payload, _ := json.Marshal(delivery)

	// resp, err := http.Post(url, "application/json",
	// 	bytes.NewBuffer(payload))
	// if err != nil || resp.StatusCode != 200 {
	// 	return fmt.Errorf("Request saga delivery fail: %w", err), resp.StatusCode
	// }

	payload, _ := json.Marshal(delivery)
	statuscode, _, err := goglib.HttpRequest(ctx, payload, "POST", url)
	if err != nil {
		logger.Log.Print(2, "Request saga delivery fail... ")
		return statuscode, fmt.Errorf("Request saga delivery fail: %w", err)
	}

	return statuscode, nil
}
