package api

import "time"

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// type orderInfoRequest struct {
// 	Username string `uri:"username" binding:"required"`
// }

type deliveryInfoRequest struct {
	Username string `form:"username" binding:"required"`
}

type deliveryRequest struct {
	OrderId int    `json:"order_id" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type orderCancel struct {
	OrderID int `uri:"order_id" binding:"required"`
}

type DeliveryResponse struct {
	DeliveryId int       `json:"id"`
	OrderId    int       `json:"order_id"`
	State      int       `json:"state"`
	Address    string    `json:"address"`
	ReqDt      time.Time `json:"req_dt"`
	ComplDt    time.Time `json:"compl_dt"`
}

type orderRequest struct {
	OrderId int `json:"order_id" binding:"required"`
	// Username   string `json:"username" binding:"required,alphanum"`
	// TotalAmout int `json:"amount" binding:"required"`
}
