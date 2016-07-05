package mp

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
)

type MenuItemType string

const (
	MenuTypePopMenu         MenuItemType = "menu"
	MenuTypeClick           MenuItemType = "click"
	MenuTypeView            MenuItemType = "view"
	MenuTypeScanCodePush    MenuItemType = "scancode_push"
	MenuTypeScanCodeWaitMsg MenuItemType = "scancode_waitmsg"
	MenuTypePicSysPhoto     MenuItemType = "pic_sysphoto"
	MenuTypePicPhotoOrAlbum MenuItemType = "pic_photo_or_album"
	MenuTypePicWeixin       MenuItemType = "pic_weixin"
	MenuTypeLocation        MenuItemType = "location_select"
	MenuTypeMedia           MenuItemType = "media_id"
	MenuTypeViewLimited     MenuItemType = "view_limited"
)

type MenuItem struct {
	Type       MenuItemType `json:"type"`
	Name       string       `json:"name"`
	Key        string       `json:"key"`
	Url        string       `json:"url"`
	MediaId    string       `json:"media_id"`
	SubButtons []MenuItem   `json:"sub_button"`
}

type Menu struct {
	Items []MenuItem `json:"button"`
}

//the struct of response from weixin server
type ResponseMenu struct {
	Menu Menu `json:"menu"`
}

//for custom json marshal
func (this *MenuItem) MarshalJSON() ([]byte, error) {
	dic := make(map[string]interface{})
	if this.Type != MenuTypePopMenu {
		dic["type"] = this.Type
	}
	if this.Type == MenuTypeView {
		dic["url"] = this.Url
	}
	if this.Type != MenuTypePopMenu && this.Type != MenuTypeView &&
		this.Type != MenuTypeMedia && this.Type != MenuTypeViewLimited {
		dic["key"] = this.Key
	}
	if this.Type == MenuTypeMedia || this.Type == MenuTypeViewLimited {
		dic["media_id"] = this.MediaId
	}
	dic["name"] = this.Name
	if this.SubButtons != nil {
		dic["sub_button"] = this.SubButtons
	}
	return json.Marshal(dic)
}

//for custom json unmarshal
func (this *MenuItem) UnmarshalJSON(data []byte) error {
	dic := make(map[string]interface{})
	err := json.Unmarshal(data, &dic)
	if err != nil {
		return err
	}
	t, ok := dic["type"]
	if ok {
		if ts, tok := t.(string); tok {
			this.Type = MenuItemType(ts)
		} else {
			return fmt.Errorf("menu's type unkown")
		}
	} else {
		this.Type = MenuTypePopMenu
	}
	if n, ok := dic["name"]; ok {
		if ts, tok := n.(string); tok {
			this.Name = ts
		}
	} else {
		return fmt.Errorf("menu's name unkown")
	}
	if n, ok := dic["key"]; ok {
		if ts, tok := n.(string); tok {
			this.Key = ts
		}
	}
	if n, ok := dic["url"]; ok {
		if ts, tok := n.(string); tok {
			this.Url = ts
		}
	}
	if n, ok := dic["media_id"]; ok {
		if ts, tok := n.(string); tok {
			this.MediaId = ts
		}
	}
	if a, ok := dic["sub_button"]; ok {
		adat, err := json.Marshal(a)
		if err != nil {
			return err
		}
		items := make([]MenuItem, 0)
		err = json.Unmarshal(adat, &items)
		if err == nil {
			this.SubButtons = items
		}
	}
	return nil
}

func (this *MenuItem) Same(another *MenuItem) bool {
	if another == nil {
		return false
	}

	if this.Type != another.Type || this.Name != another.Name {
		return false
	}
	switch this.Type {
	case MenuTypePopMenu:
		a, b := 0, 0
		if this.SubButtons != nil {
			a = len(this.SubButtons)
		}
		if another != nil {
			b = len(another.SubButtons)
		}
		if a != b {
			return false
		}
		if a > 0 {
			for i, m := range this.SubButtons {
				if !m.Same(&another.SubButtons[i]) {
					return false
				}
			}
		}
	case MenuTypeView:
		return this.Url == another.Url
	case MenuTypeMedia, MenuTypeViewLimited:
		return this.MediaId == another.MediaId
	default:
		return this.Key == another.Key
	}
	return true
}

func (this *Menu) Same(another *Menu) bool {
	a, b := 0, 0
	if this.Items != nil {
		a = len(this.Items)
	}
	if another.Items != nil {
		b = len(this.Items)
	}
	if a != b {
		return false
	}
	for i, mi := range this.Items {
		if !mi.Same(&another.Items[i]) {
			return false
		}
	}
	return true
}

//query menu from weixin server
func (this *WxMp) GetMenu() (menu *Menu, err error) {
	dat, err := this.HttpGet("/cgi-bin/menu/get", true, nil)
	if err != nil {
		logrus.Error(`获取菜单错误`, err)
		return nil, err
	}
	if len(dat) == 0 {
		logrus.Info(`未设置菜单`)
		return nil, nil
	}
	logrus.Debug(`存在菜单数据`)
	respMenu := ResponseMenu{}
	err = json.Unmarshal(dat, &respMenu)
	if err != nil {
		logrus.Error(`获取菜单数据错误`, err)
		return nil, err
	}
	return &respMenu.Menu, err
}

//set weixin mp's menu
func (this *WxMp) SetMenu(menu *Menu) error {
	_, err := this.HttpPost("/cgi-bin/menu/create", true, nil, menu)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"menu": menu,
			"err":  err,
		}).Error("set weixin's menu")
	}
	return err
}

//update menu,if menu is same with service menu,do nothing
func (this *WxMp) UpdateMenu(menu *Menu) error {
	oMenu, err := this.GetMenu()
	if err != nil {
		return nil
	}
	if oMenu == nil {
		return this.SetMenu(menu)
	}
	logrus.WithFields(logrus.Fields{
		"old menu": oMenu,
		"new menu": menu,
	}).Info("set weixin's menu")
	if !oMenu.Same(menu) {
		return this.SetMenu(menu)
	}
	logrus.Info("same menu,do nothing")
	//do nothing
	return nil
}

//delete menu
func (this *WxMp) DeleteMenu() error {
	_, err := this.HttpGet("/cgi-bin/menu/delete", true, nil)
	return err
}
