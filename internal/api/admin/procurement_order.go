package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/tools/a2r"
)

func (o *Api) AddProcurementOder(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.AddProcurementOder, o.adminClient)
}

func (o *Api) DeleteProcurementOder(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DeleteProcurementOder, o.adminClient)
}

func (o *Api) UpdateProcurementOder(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.UpdateProcurementOder, o.adminClient)
}

func (o *Api) GetProcurementOder(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.GetProcurementOder, o.adminClient)
}

func (o *Api) QueryProcurementOderList(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.QueryProcurementOderList, o.adminClient)
}

func (o *Api) CountProcurementOderByStatus(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.CountProcurementOderByStatus, o.adminClient)
}
