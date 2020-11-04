package api

import (
	"strconv"
	"time"

	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/lib/md5"
	"github.com/tang-go/go-dog/plugins"
)

//GetAdminInfo 获取管理员信息
func (pointer *API) GetAdminInfo(ctx plugins.Context, request param.GetAdminInfoReq) (response param.GetAdminInfoRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	response.ID = strconv.FormatInt(admin.AdminID, 10)
	response.Name = admin.Name
	response.Avatar = "/avatar2.jpg"
	response.Phone = admin.Phone
	role := new(table.OwnerRole)
	if pointer.mysql.GetReadEngine().Where("role_id = ?", admin.RoleID).First(role).RecordNotFound() == true {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员权限不正确")
		return
	}
	response.RoleID = role.Name
	response.Role.ID = role.Name
	if role.IsAdmin {
		response.Role.Permissions = append(response.Role.Permissions, &param.Permissions{
			RoleID:         role.Name,
			PermissionID:   "admin",
			PermissionName: role.Description,
			ActionEntitySet: []*param.ActionEntitySet{
				&param.ActionEntitySet{
					Action:       "add",
					Describe:     "新增",
					DefaultCheck: true,
				},
				&param.ActionEntitySet{
					Action:       "delete",
					Describe:     "删除",
					DefaultCheck: true,
				},
				&param.ActionEntitySet{
					Action:       "update",
					Describe:     "修改",
					DefaultCheck: true,
				},
				&param.ActionEntitySet{
					Action:       "query",
					Describe:     "查询",
					DefaultCheck: true,
				},
			},
		})
	}
	return
}

//AdminLogin 管理员登录
func (pointer *API) AdminLogin(ctx plugins.Context, request param.AdminLoginReq) (response param.AdminLoginRes, err error) {
	//查询是否拥有此用户
	admin := new(table.Admin)
	if pointer.mysql.GetReadEngine().Where("phone = ?", request.Phone).First(admin).RecordNotFound() == true {
		err = customerror.EnCodeError(define.AdminLoginErr, "管理员登录失败")
		return
	}
	//密码对比
	if md5.Md5(md5.Md5(request.Pwd)+admin.Salt) != admin.Pwd {
		err = customerror.EnCodeError(define.AdminLoginErr, "管理员登录失败")
		return
	}
	//生成登录记录
	mysqllog := &table.Log{
		//日志ID
		LogID: pointer.snowflake.GetID(),
		//类型
		Type: table.LoginType,
		//操作人
		AdminID: admin.AdminID,
		//名称
		AdminName: admin.Name,
		//操作方法
		Method: "AdminLogin",
		//描述
		Description: "管理员登录",
		//业主ID
		OwnerID: admin.OwnerID,
		//操作IP
		IP: ctx.GetAddress(),
		//操作URL
		URL: ctx.GetDataByKey("URL").(string),
		//操作时间
		Time: time.Now().Unix(),
	}
	if e := pointer.mysql.GetWriteEngine().Create(mysqllog).Error; e != nil {
		err = customerror.EnCodeError(define.AdminLoginErr, e.Error())
		return
	}
	//生成token
	token := md5.Md5(admin.AdminID)
	//生成token缓存
	pointer.cache.GetCache().SetByTime(token, admin, define.AdminTokenValidityTime)
	//登录成功返回
	response.Name = admin.Name
	response.OwnerID = strconv.FormatInt(admin.OwnerID, 10)
	response.Token = token
	return
}