package router

import (
	"webOrder/service"

	"github.com/gin-gonic/gin"
)

type Router struct {
	RouterServer *service.OrderService
}

func NewEngine(h *Router) *gin.Engine {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	//resfulapi
	r.GET("/order/:order_no", h.QueryOrderByNo())
	r.POST("/order/", h.NewOrder())
	r.PUT("/order/", h.UpdateOrderByOrderNo())
	r.GET("/order/like/:key_word", h.QueryLikeOrder())
	r.POST("/upload/", h.Upload())
	r.GET("/order/all/", h.GetAllOrder())
	r.GET("/order/excel/", h.GetExclfile())
	r.GET("/order/list/:key_word", h.QueryOrderList())
	r.GET("/order/file/:order_no", h.Download())

	return r
}