package mp

import (
	"encoding/json"
	"testing"
)

func TestMenuMarshalJSON(t *testing.T) {
	return
	menu := &Menu{
		Items: make([]MenuItem, 0),
	}
	menu.Items = append(menu.Items, MenuItem{
		Type: MenuTypeClick,
		Name: "今日歌曲",
		Key:  "V1001_TODAY_MUSIC",
	})
	menuItem := MenuItem{
		Type:       MenuTypePopMenu,
		Name:       "菜单",
		SubButtons: make([]MenuItem, 0),
	}
	menuItem.SubButtons = append(menuItem.SubButtons, MenuItem{
		Type: MenuTypeView,
		Name: "搜索",
		Url:  "http://www.soso.com/",
	})
	menuItem.SubButtons = append(menuItem.SubButtons, MenuItem{
		Type: MenuTypeView,
		Name: "视频",
		Url:  "http://v.qq.com/",
	})
	menuItem.SubButtons = append(menuItem.SubButtons, MenuItem{
		Type: MenuTypeClick,
		Name: "赞一下我们",
		Key:  "V1001_GOOD",
	})
	menu.Items = append(menu.Items, menuItem)
	dat, err := json.Marshal(menu)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%v", string(dat))
	}
}

func TestMenuUnmarshalJSON(t *testing.T) {
	return
	str := `{"button":[{"key":"V1001_TODAY_MUSIC","name":"今日歌曲","type":"click"},{"name":"菜单","sub_button":[{"name":"搜索","type":"view","url":"http://www.soso.com/"},{"name":"视频","type":"view","url":"http://v.qq.com/"},{"key":"V1001_GOOD","name":"赞一下我们","type":"click"}]}]}`
	menu := Menu{}
	err := json.Unmarshal([]byte(str), &menu)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%v", menu)
	}
}

func TestMenuSame(t *testing.T) {
	return
	wx := NewWxMp(config)
	menu, err := wx.GetMenu()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", menu)

	str := `{"button":[{"name":"点我","sub_button":[{"type":"click","name":"电子会员卡","key":"vip"},{"type":"view","name":"会员特权","url":"http://weixintest.imooly.com/show/"},{"type":"view","name":"充值卡","url":"http://weixintest.imooly.com/account/card"}]},{"name":"随便点","sub_button":[{"type":"click","name":"热门店铺","key":"HOTBUSINESS"},{"type":"click","name":"附近店铺","key":"NEARBUSINESS"}]},{"name":"不要点我","sub_button":[{"type":"view","name":"常见问题","url":"http://weixintest.imooly.com/show/problem"},{"type":"view","name":"调戏客服","url":"http://weixintest.imooly.com/show/contact"}]}]}`
	menub := Menu{}
	err = json.Unmarshal([]byte(str), &menub)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", menub)
	if menu.Same(&menub) {
		t.Log("相同")
	} else {
		t.Log("不同")
	}
}
