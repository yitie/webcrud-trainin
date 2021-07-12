package router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx/v3"
	"github.com/w3liu/go-common/constant/timeformat"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
	"webOrder/model"
	fmtdata "webOrder/server/server_model"
)

var num int64

func Generate(t time.Time) string {
	s := t.Format(timeformat.Continuity)
	m := t.UnixNano()/1e6 - t.UnixNano()/1e9*1e3
	ms := sup(m, 3)
	p := os.Getpid() % 1000
	ps := sup(int64(p), 3)
	i := atomic.AddInt64(&num, 1)
	r := i % 10000
	rs := sup(r, 4)
	n := fmt.Sprintf("%s%s%s%s", s, ms, ps, rs)
	return n
}
func sup(i int64, n int) string {
	m := fmt.Sprintf("%d", i)
	for len(m) < n {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}

func (h *Router) httpErr(c *gin.Context, code int, err interface{}) {
	c.JSON(code, gin.H{
		"status": "Failure",
		"err":    err,
	})
}

func (h *Router) httpSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"err":    data,
	})
}

func (h *Router) NewOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req *fmtdata.AddOrderReq
		if err := c.BindJSON(&req); err != nil {
			h.httpErr(c, http.StatusBadRequest, err.Error())
			return
		}

		if err := req.IsValid(); err != nil {
			h.httpErr(c, http.StatusBadRequest, err.Error())
			return
		}

		o := model.Order{
			OrderNo:  fmt.Sprintf("%s", Generate(time.Now())),
			UserName: req.UserName,
			Amount:   req.Amount,
			Status:   true,
			FileUrl:  req.FileUrl,
		}

		if err := h.RouterServer.NewOrder(&o); err != nil {
			h.httpErr(c, http.StatusBadGateway, err.Error())
			return
		}

		h.httpSuccess(c, "success create data ")

	}
}

func (h *Router) QueryOrderByNo() gin.HandlerFunc {
	return func(c *gin.Context) {
		no := c.Param("order_no")
		if no == "" {
			errors.New("U query without no")
			h.httpErr(c, http.StatusBadGateway, "update without no")
			return
		}
		order, err := h.RouterServer.QueryOrderByNo(no)
		if err != nil {
			h.httpErr(c, http.StatusBadGateway, err.Error())
			return
		}
		if order == nil {
			h.httpErr(c, http.StatusBadRequest, "未查询到订单")
			return
		}

		h.httpSuccess(c, &model.Order{
			Status:   order.Status,
			OrderNo:  order.OrderNo,
			UserName: order.UserName,
			Amount:   order.Amount,
			FileUrl:  order.FileUrl,
		})
	}
}

func (h *Router) UpdateOrderByOrderNo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req *fmtdata.UpdateOrder

		if err := req.IsValid(); err != nil {
			h.httpErr(c, http.StatusBadRequest, err.Error())
			return
		}
		no := c.PostForm("order_no")
		status := c.PostForm("status")
		fileUrl := c.PostForm("file_url")
		amount := c.PostForm("amount")
		tempAmount, _ := strconv.ParseFloat(amount, 64)

		err := h.RouterServer.UpdateOrderByOrderNo(no, map[string]interface{}{
			"Amount":  tempAmount, //amount
			"Status":  status,
			"FileUrl": fileUrl,
		})

		if err != nil {
			h.httpErr(c, http.StatusBadGateway, err.Error())
			return
		}
		h.httpSuccess(c, nil)

	}
}

func (h *Router) GetAllOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var orders []*model.Order
		err := h.RouterServer.GetAllOrder(orders)
		if err != nil {
			h.httpErr(c, http.StatusBadRequest, err.Error())
			return
		}
		h.httpSuccess(c, orders)
	}
}

func (h *Router) QueryLikeOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		keyWord := c.Param("key_word")

		if keyWord == "" {
			h.httpErr(c, http.StatusBadGateway, "U query without no")
			return
		}

		orders, err := h.RouterServer.QueryLikeOrder(keyWord)
		if err != nil {
			h.httpErr(c, http.StatusBadRequest, err.Error())
			return
		}
		if orders == nil {
			h.httpErr(c, http.StatusBadRequest, "未查询到结果")
			return
		}

		h.httpSuccess(c, orders)
	}
}

func (h *Router) QueryOrderList() gin.HandlerFunc {
	return func(c *gin.Context) {
		keyWord := c.Param("key_word")

		if keyWord == "" {
			h.httpErr(c, http.StatusBadGateway, "U query without no")
			return
		}

		up := false
		orders, err := h.RouterServer.QueryOrderList(keyWord, up)
		if err != nil {
			h.httpErr(c, http.StatusBadRequest, err.Error())
			return
		}
		if orders == nil {
			h.httpErr(c, http.StatusBadRequest, "未查询到结果")
			return
		}
		h.httpSuccess(c, orders)

	}
}

func (h *Router) Upload() gin.HandlerFunc {
	return func(c *gin.Context) {
		no := c.PostForm("order_no")

		// 单文件
		file, err := c.FormFile("file")//想处理但是file没获得是nil吗？
		if err != nil {
			c.JSON(403, gin.H{
				"error": err.Error(),
			})
			return

		} else {
			dst := fmt.Sprintf("./static/%s", file.Filename)

			// 上传文件至指定目录
			c.SaveUploadedFile(file, dst)

			err := h.RouterServer.UpdateOrderByOrderNo(no, map[string]interface{}{
				"FileUrl": dst,
			})
			if err != nil {
				h.httpErr(c, http.StatusBadRequest, err.Error())
				return
			}
			h.httpSuccess(c, nil)
		}
	}
}

func (h *Router) Download() gin.HandlerFunc {
	return func(c *gin.Context) {
		no := c.Param("order_no")
		order, err := h.RouterServer.QueryOrderByNo(no)
		if err != nil {
			h.httpErr(c, http.StatusBadRequest, err.Error())
			return
		}
		if order == nil {
			h.httpErr(c, http.StatusBadRequest, "未查询到订单")
			return
		}
		fmt.Println(order.FileUrl)
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", order.FileUrl)) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
		c.File(order.FileUrl)
	}
}

func (h *Router) GetExclfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var orders []*model.Order
		err := h.RouterServer.GetAllOrder(orders)
		if err != nil {
			h.httpErr(c, http.StatusBadRequest, err.Error())
			return
		}
		dst := "./static/export.xlsx"
		wb := xlsx.NewFile()
		sheet, err := wb.AddSheet("demo_order_list")
		if err != nil {
			fmt.Printf(err.Error())
			return
		}

		for _, d := range orders {
			row := sheet.AddRow()
			ID := row.AddCell()
			ID.SetValue(d.ID)

			OrderNo := row.AddCell()
			OrderNo.SetValue(d.OrderNo)

			UserName := row.AddCell()
			UserName.SetValue(d.UserName)

			Amount := row.AddCell()
			Amount.SetValue(d.Amount)

			Status := row.AddCell()
			Status.SetValue(d.Status)

			FileUrl := row.AddCell()
			FileUrl.SetValue(d.Status)
		}
		err = wb.Save(dst)
		if err != nil {
			panic("save file err")
		}
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", dst)) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
		c.File(dst)
	}

}