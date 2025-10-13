package api

import "delivery_service/internal/db"

func convertDelivery(deli db.DELIVERIES) DeliveryResponse {

	return DeliveryResponse{
		DeliveryId: deli.DELIVERY_ID,
		OrderId:    deli.ORDER_ID,
		State:      deli.STATUS,
		Address:    deli.ADDRESS,
		ReqDt:      deli.REQ_DT,
		ComplDt:    deli.COMPL_DT.Time,
	}
}

func getDeliveryPrarm(deli deliveryRequest) db.DELIVERIES {
	return db.DELIVERIES{
		ORDER_ID: deli.OrderId,
		ADDRESS:  deli.Address,
	}
}
