package server_model

import "errors"

type AddOrderReq struct {
	UserName string  `gorm:"index:idx_order_username" json:"user_name"`
	OrderNo  string  `gorm:"index:idx_order_no" json:"order_no"`
	Amount   float64 `json:"amount"`
	FileUrl  string  `json:"file_url"`
}

func (r *AddOrderReq) IsValid() error {
	switch {
	case r.Amount == 0:
		return errors.New("Amount is not zero. ")
	case r.UserName == "":
		return errors.New("UserName is not found. ")

	default:
		return nil
	}
}

type UpdateOrder struct {
	ID      uint    `gorm:"primarykey"`
	OrderNo string  `gorm:"index:idx_order_no" json:"order_no"`
	Amount  float64 `json:"amount"`
	FileUrl string  `json:"file_url"`
	Status  string  `json:"status"`
}

func (r *UpdateOrder) IsValid() error {
	switch {
	case r.OrderNo == "":
		return errors.New("OrderNo is not found. ")
	case r.Amount == 0:
		return errors.New("Amount is not zero")
	case r.Status == "":
		return errors.New("Status is not found")
	default:
		return nil
	}
}