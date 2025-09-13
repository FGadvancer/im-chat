package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/tools/a2r"
)

func (o *Api) AddRebateOder(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.AddRebateOder, o.adminClient)
}

func (o *Api) DeleteRebateOder(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DeleteRebateOder, o.adminClient)
}

func (o *Api) UpdateRebateOder(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.UpdateRebateOder, o.adminClient)
}

func (o *Api) GetRebateOder(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.GetRebateOder, o.adminClient)
}

func (o *Api) QueryRebateOderList(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.QueryRebateOderList, o.adminClient)
}

func (o *Api) CountRebateOderByStatus(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.CountRebateOderByStatus, o.adminClient)
}
