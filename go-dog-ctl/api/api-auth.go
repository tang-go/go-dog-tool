package api

import (
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/plugins"
)

//Auth 插件
func (pointer *API) Auth(ctx plugins.Context, method, token string) error {
	admin := new(table.Admin)
	if e := pointer.cache.GetCache().Get(token, admin); e != nil {
		return customerror.EnCodeError(define.AdminTokenErr, "token失效或者不正确")
	}
	ctx.SetShare("Admin", admin)
	return nil
}
