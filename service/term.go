package service

import (
	"fmt"
	"goblog/dao"
	"goblog/model"
	"goblog/model/request"
)

// 编辑项，返回该项ID（新建为新的项ID，修改则与输入相同）
func EditTerm(req *request.Term) (uint32, error) {
	dao := dao.NewDao()
	var v bool
	// 校验项类型
	termType := model.TermType(req.Type)
	if termType != model.TERM_TYPE_CATEGORY && termType != model.TERM_TYPE_TAG {
		return 0, fmt.Errorf("项类型非法")
	}

	// 构造项
	var term *model.Term
	if req.ID == 0 {
		// 新建前，slug必须不存在
		_, v = dao.Term().GetBySlug(req.Slug)
		if v {
			return 0, fmt.Errorf("slug对应的项已存在")
		}
		// 构造新的项
		term = &model.Term{
			Tid:   0,
			Slug:  req.Slug,
			Type:  termType,
			Count: 0,
		}
	} else {
		// 修改前，ID对应的项必须存在
		term, v = dao.Term().GetByID(req.ID)
		if !v {
			return 0, fmt.Errorf("ID对应的项不存在")
		}
	}
	term.Name = req.Name
	term.Description = req.Description

	// 提交数据库
	_ = dao.Begin()
	if term.Tid == 0 {
		term.Tid, v = dao.Term().Add(term)
		if !v {
			_ = dao.Rollback()
			return 0, fmt.Errorf("新建项失败")
		}
	} else {
		v = dao.Term().Modify(term)
		if !v {
			_ = dao.Rollback()
			return 0, fmt.Errorf("修改项失败")
		}
	}
	err := dao.Commit()
	if err != nil {
		_ = dao.Rollback()
		return 0, fmt.Errorf("提交事务失败")
	}
	return term.Tid, nil
}

// 获取单个项
func GetTerm(id uint32) (*model.Term, error) {
	dao := dao.NewDao()
	// 获取项
	term, v := dao.Term().GetByID(id)
	if !v {
		return nil, fmt.Errorf("项不存在")
	}
	return term, nil
}

// 获取项列表
func ListTerm(termType uint8, pi, ps uint32) ([]model.Term, error) {
	dao := dao.NewDao()
	// 校验项类型
	t := model.TermType(termType)
	if t != model.TERM_TYPE_CATEGORY && t != model.TERM_TYPE_TAG {
		return nil, fmt.Errorf("项类型非法")
	}

	// 获取项列表
	list, v := dao.Term().ListByType(t, true, pi, ps)
	if !v {
		return nil, fmt.Errorf("该页不存在项")
	}
	return list, nil
}

// 统计项数量
func CountTermByType(termType uint8) (uint32, error) {
	// 校验项类型
	t := model.TermType(termType)
	if t != model.TERM_TYPE_CATEGORY && t != model.TERM_TYPE_TAG {
		return 0, fmt.Errorf("项类型非法")
	}

	return dao.NewDao().Term().CountByType(t), nil
}
