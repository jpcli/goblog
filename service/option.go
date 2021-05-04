package service

import (
	"fmt"
	"goblog/dao"
	"goblog/model"
	"goblog/model/request"
)

// 设置选项(泛用)，考虑到程序设计，添加选项暂不执行
func SetOption(kvs map[string]string) error {
	dao := dao.NewDao()
	// 判断新增还是修改
	addGroup, modifyGroup := make(map[string]string), make(map[string]string)
	for name, val := range kvs {
		optVal, v := dao.Opt().GetByName(name)
		if v {
			if optVal != val {
				modifyGroup[name] = val
			}
		} else {
			addGroup[name] = val
		}
	}

	// 开始数据库操作
	_ = dao.Begin()
	// TODO 添加选项

	// 修改选项
	for name, val := range modifyGroup {
		v := dao.Opt().Modify(&model.Option{Name: name, Value: val})
		if !v {
			_ = dao.Rollback()
			return fmt.Errorf("修改选项失败")
		}
	}
	err := dao.Commit()
	if err != nil {
		_ = dao.Rollback()
		return fmt.Errorf("提交事务失败")
	}
	return nil
}

// 获取选项(泛用),结果返回到参数中的map里
func GetOption(kvs map[string]string) error {
	dao := dao.NewDao()
	for name := range kvs {
		val, v := dao.Opt().GetByName(name)
		if !v {
			return fmt.Errorf("[%s]选项不存在", name)
		} else {
			kvs[name] = val
		}
	}
	return nil
}

// 设置基本选项
func SetBaseOption(req *request.BaseOption) error {
	m := make(map[string]string)
	m["page_size"] = fmt.Sprintf("%d", req.PageSize)
	m["page_nav_size"] = fmt.Sprintf("%d", req.PageNavSize)
	m["site_name"] = req.SiteName
	m["site_url"] = req.SiteURL
	return SetOption(m)
}

// 获取基本选项
func GetBaseOption() (map[string]string, error) {
	m := make(map[string]string)
	m["page_size"] = ""
	m["page_nav_size"] = ""
	m["site_name"] = ""
	m["site_url"] = ""
	err := GetOption(m)
	return m, err
}
