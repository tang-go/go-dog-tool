package service

import (
	"sort"
	"time"

	"github.com/tang-go/go-dog-tool/define"
	authAPI "github.com/tang-go/go-dog-tool/go-dog-auth/api"
	authTable "github.com/tang-go/go-dog-tool/go-dog-auth/table"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//MenuSort 菜单排序
type MenuSort []*param.Menu

//Len 长度
func (m MenuSort) Len() int { return len(m) }

//Swap 交换
func (m MenuSort) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

//Less 比较
func (m MenuSort) Less(i, j int) bool {
	if m[i].Sort == m[j].Sort {
		return m[i].Time > m[j].Time
	}
	return m[i].Sort > m[j].Sort
}

//_AssembleMenu 组织菜单
func (s *Service) _AssembleMenu(m *param.Menu, source []authTable.SysMenu, parentID uint) {
	for _, menu := range source {
		if menu.ParentID == parentID {
			nm := param.Menu{
				ID:       menu.ID,
				ParentID: menu.ParentID,
				Organize: menu.Organize,
				URL:      menu.URL,
				Describe: menu.Describe,
				Sort:     menu.Sort,
				Time:     time.Unix(menu.Time, 0).Format("2006-01-02 15:04:05"),
			}
			s._AssembleMenu(&nm, source, menu.ID)
			m.Children = append(m.Children, &nm)
		}
	}
	sort.Sort(MenuSort(m.Children))
	return
}

//GetMenu 获取菜单
func (s *Service) GetMenu(ctx plugins.Context, request param.GetMenuReq) (response param.GetMenuRes, err error) {
	//获取菜单信息
	menus, e := authAPI.SelectMenu(ctx, define.Organize)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员角色菜单失败")
		return
	}
	for _, menu := range menus {
		if menu.ParentID == 0 {
			nm := param.Menu{
				ID:       menu.ID,
				ParentID: menu.ParentID,
				Organize: menu.Organize,
				URL:      menu.URL,
				Describe: menu.Describe,
				Sort:     menu.Sort,
				Time:     time.Unix(menu.Time, 0).Format("2006-01-02 15:04:05"),
			}
			s._AssembleMenu(&nm, menus, menu.ID)
			response.Menu = append(response.Menu, &nm)
		}
	}
	sort.Sort(MenuSort(response.Menu))
	return
}

//CreateMenu 创建菜单
func (s *Service) CreateMenu(ctx plugins.Context, request param.CreateMenuReq) (response param.CreateMenuRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetBuildServiceListErr, "管理员信息失败")
		return
	}
	if _, e := authAPI.CreateMenu(ctx, define.Organize, request.Describe, request.URL, request.ParentID, request.Sort); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.CreateMenuErr, "创建菜单失败")
		return
	}
	//生成登录记录
	tbLog := &table.Log{
		LogID:       s.snowflake.GetID(),
		Type:        table.CreateMenuType,
		AdminID:     admin.AdminID,
		AdminName:   admin.Name,
		Method:      "CreateMenu",
		Description: "创建菜单",
		OwnerID:     admin.OwnerID,
		IP:          ctx.GetAddress(),
		URL:         ctx.GetDataByKey("URL").(string),
		Time:        time.Now().Unix(),
	}
	if e := s.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
	}
	response.Success = true
	return
}
