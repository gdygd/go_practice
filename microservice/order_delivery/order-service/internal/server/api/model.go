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

type orderInfoRequest struct {
	Username string `form:"username" binding:"required"`
}

type orderRequest struct {
	Username   string `json:"username" binding:"required,alphanum"`
	TotalAmout int    `json:"amount" binding:"required"`
}

type orderCancel struct {
	OrderID int `uri:"order_id" binding:"required"`
}

type OrderResponse struct {
	OrderId   int       `json:"id"`
	Username  string    `json:"username"`
	State     int       `json:"state"`
	OrderDT   time.Time `json:"order_dt"`
	TotAmount int       `json:"tot_amount"`
}
