package api

import "order-service/internal/db"

func convertOrder(order db.ORDER) OrderResponse {

	return OrderResponse{
		OrderId:   order.ORDER_ID,
		Username:  order.USER_NM,
		State:     order.STATE,
		OrderDT:   order.ORDER_DT,
		TotAmount: order.TOT_AMOUNT,
	}
}

func getOrderPrarm(order orderRequest) db.ORDER {
	return db.ORDER{
		USER_NM:    order.Username,
		TOT_AMOUNT: order.TotalAmout,
	}
}
