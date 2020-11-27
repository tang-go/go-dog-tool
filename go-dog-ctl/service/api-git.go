package service

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//CreateHarbor 创建git账号
func (s *Service) CreateGit(ctx plugins.Context, request param.CreateGitReq) (response param.CreateGitRes, err error) {
	if request.Address == "" {
		err = customerror.EnCodeError(define.CreateGitErr, "git仓库地址不能为空")
		return
	}
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	git := new(table.Git)
	if s.mysql.GetReadEngine().Where("owner_id = ? AND address = ?", admin.OwnerID, request.Address).First(git).RecordNotFound() == true {
		if e := s._CreateLog(admin, table.CreateGitType, ctx.GetAddress(), ctx.GetURL(), "CreateGit", "创建git账号", func(tx *gorm.DB) error {
			if e := tx.Create(&table.Git{
				Address: request.Address,
				Account: request.Account,
				Pwd:     request.Pwd,
				OwnerID: admin.OwnerID,
				Time:    time.Now().Unix(),
			}).Error; e != nil {
				log.Errorln(e.Error())
				return e
			}
			return nil
		}); e != nil {
			log.Errorln(e.Error())
			err = customerror.EnCodeError(define.CreateGitErr, "创建git账号失败")
			return
		}
		response.Success = true
		return
	}
	err = customerror.EnCodeError(define.CreateGitErr, "创建git账号失败")
	return
}

//DelGit 删除git账号
func (s *Service) DelGit(ctx plugins.Context, request param.DelGitReq) (response param.DelGitRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	if e := s._CreateLog(admin, table.DelGitType, ctx.GetAddress(), ctx.GetURL(), "DelGit", "删除git账号", func(tx *gorm.DB) error {
		if e := tx.Where("owner_id = ? AND id = ?", admin.OwnerID, request.ID).Delete(&table.Git{}).Error; e != nil {
			log.Errorln(e.Error())
			return e
		}
		return nil
	}); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.DelGitErr, "删除git账号失败")
		return
	}
	response.Success = true
	return
}

//GetGitList 获取gits账号列表
func (s *Service) GetGitList(ctx plugins.Context, request param.GetGitListReq) (response param.GetGitListRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	var gits []table.Git
	e := s.mysql.GetReadEngine().Where("owner_id = ?", admin.OwnerID).Find(&gits).Error
	if e != nil {
		err = customerror.EnCodeError(define.GetGitListErr, "获取Git列表失败")
		return
	}
	for _, git := range gits {
		response.Gits = append(response.Gits, param.Git{
			ID:      git.ID,
			Address: git.Address,
			Account: git.Account,
			Pwd:     git.Pwd,
			Time:    time.Unix(git.Time, 0).Format("2006-01-02 15:04:05"),
		})
	}
	return
}
