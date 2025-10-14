package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func callOrderServiceConfirm(orderId int) error {

	return nil
}

func callOrderServiceCancel(orderId int) error {

	return nil
}

func callDeliveryService(orderId int) error {
	var delivery deliveryRequest
	delivery.OrderId = orderId
	delivery.Address = "gunpo1"

	payload, _ := json.Marshal(delivery)
	url := "http://10.1.0.119:9080/delivery/request"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("delivery 실패: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
