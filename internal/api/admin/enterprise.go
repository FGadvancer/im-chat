package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/tools/a2r"
)

func (o *Api) AddEnterpriseInfo(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.AddEnterpriseInfo, o.adminClient)
}

func (o *Api) DeleteEnterpriseInfo(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.DeleteEnterpriseInfo, o.adminClient)
}

func (o *Api) UpdateEnterpriseInfo(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.UpdateEnterpriseInfo, o.adminClient)
}

func (o *Api) GetEnterpriseInfo(c *gin.Context) {
	a2r.Call(c, admin.AdminClient.GetEnterpriseInfo, o.adminClient)
}
