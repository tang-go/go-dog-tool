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

//CreateDocs 创建文档
func (s *Service) CreateDocs(ctx plugins.Context, request param.CreateDocsReq) (response param.CreateDocsRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	if e := s.mysql.GetReadEngine().Where("url = ? AND owner_id = ?", request.URL, admin.OwnerID).First(new(table.Docs)).Error; e != nil {
		if e.Error() == gorm.ErrRecordNotFound.Error() {
			if e := s._CreateLog(admin, table.CreateDocsType, ctx.GetAddress(), ctx.GetURL(), "CreateDocs", "创建文档", func(tx *gorm.DB) error {
				if e := tx.Create(&table.Docs{
					OwnerID: admin.OwnerID,
					Name:    request.Name,
					URL:     request.URL,
					Time:    time.Now().Unix(),
				}).Error; e != nil {
					log.Errorln(e.Error())
					return e
				}
				return nil
			}); e != nil {
				log.Errorln(e.Error())
				err = customerror.EnCodeError(define.CreateDocsErr, "创建文档失败")
				return
			}
			response.Success = true
			return
		}
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.CreateDocsErr, e.Error())
		return
	}
	err = customerror.EnCodeError(define.CreateDocsErr, "存在此url地址的文档")
	return
}

//DelDocs 删除文档
func (s *Service) DelDocs(ctx plugins.Context, request param.DelDocsReq) (response param.DelDocsRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	if e := s._CreateLog(admin, table.DelDocsType, ctx.GetAddress(), ctx.GetURL(), "DelDocs", "删除文档", func(tx *gorm.DB) error {
		if e := tx.Where("owner_id = ? AND id = ?", admin.OwnerID, request.ID).Delete(&table.Docs{}).Error; e != nil {
			log.Errorln(e.Error())
			return e
		}
		return nil
	}); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.DelDocsErr, "删除文档失败")
		return
	}
	response.Success = true
	return
}

//GetDocs 获取文档
func (s *Service) GetDocs(ctx plugins.Context, request param.GetDocsReq) (response param.GetDocsRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	if e := s.mysql.GetReadEngine().Where("owner_id = ?", admin.OwnerID).Find(&response.Docs).Error; e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	return
}
