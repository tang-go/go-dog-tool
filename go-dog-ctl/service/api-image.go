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

//CreateHarbor 创建image账号
func (s *Service) CreateImage(ctx plugins.Context, request param.CreateImageReq) (response param.CreateImageRes, err error) {
	if request.Address == "" {
		err = customerror.EnCodeError(define.CreateImageErr, "image仓库地址不能为空")
		return
	}
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	image := new(table.Image)
	if s.mysql.GetReadEngine().Where("owner_id = ? AND address = ?", admin.OwnerID, request.Address).First(image).RecordNotFound() == true {
		if e := s._CreateLog(admin, table.CreateImageType, ctx.GetAddress(), ctx.GetURL(), "CreateImage", "创建image账号", func(tx *gorm.DB) error {
			if e := tx.Create(&table.Image{
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
			err = customerror.EnCodeError(define.CreateImageErr, "创建image账号失败")
			return
		}
		response.Success = true
		return
	}
	err = customerror.EnCodeError(define.CreateImageErr, "创建image账号失败")
	return
}

//DelImage 删除image账号
func (s *Service) DelImage(ctx plugins.Context, request param.DelImageReq) (response param.DelImageRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	if e := s._CreateLog(admin, table.DelImageType, ctx.GetAddress(), ctx.GetURL(), "DelImage", "删除image账号", func(tx *gorm.DB) error {
		if e := tx.Where("owner_id = ? AND id = ?", admin.OwnerID, request.ID).Delete(&table.Image{}).Error; e != nil {
			log.Errorln(e.Error())
			return e
		}
		return nil
	}); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.DelImageErr, "删除image账号失败")
		return
	}
	response.Success = true
	return
}

//GetImageList 获取images账号列表
func (s *Service) GetImageList(ctx plugins.Context, request param.GetImageListReq) (response param.GetImageListRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	var images []table.Image
	e := s.mysql.GetReadEngine().Where("owner_id = ?", admin.OwnerID).Find(&images).Error
	if e != nil {
		err = customerror.EnCodeError(define.GetImageListErr, "获取镜像列表失败")
		return
	}
	for _, image := range images {
		response.Images = append(response.Images, param.Image{
			ID:      image.ID,
			Address: image.Address,
			Account: image.Account,
			Pwd:     image.Pwd,
			Time:    time.Unix(image.Time, 0).Format("2006-01-02 15:04:05"),
		})
	}
	return
}
