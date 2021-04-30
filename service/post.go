package service

import (
	"fmt"
	"goblog/dao"
	"goblog/model"
	"goblog/model/request"
	"time"

	"github.com/spf13/cast"
)

// 修改文章，返回该文章ID（新建则为新文章ID，否则与输入一致）
func EditPost(req *request.Post) (uint32, error) {
	dao := dao.NewDao()
	var v bool

	// 校验分类
	cate, v := dao.Term().GetByID(req.CateID)
	if !v || cate.Type != model.TERM_TYPE_CATEGORY {
		return 0, fmt.Errorf("文章分类不存在")
	}

	// 校验标签
	for _, tagID := range req.TagsID {
		tag, v := dao.Term().GetByID(tagID)
		if !v || tag.Type != model.TERM_TYPE_TAG {
			return 0, fmt.Errorf("文章中有不存在的标签")
		}
	}

	// 构造文章
	var post *model.Post
	now := cast.ToUint32(time.Now().Unix())
	if req.ID > 0 {
		var has bool
		if post, has = dao.Post().GetByID(req.ID); !has {
			return 0, fmt.Errorf("文章ID不存在")
		}
	} else {
		post = &model.Post{
			Created:      now,
			CommentCount: 0,
		}
	}
	post.Modified = now
	post.Title = req.Title
	post.Keywords = req.Keywords
	post.Text = req.Text
	if a := model.PostCommentAllowType(req.CommentAllow); a == model.POST_COMMENT_ALLOW || a == model.POST_COMMENT_DISALLOW {
		post.CommentAllow = a
	} else {
		return 0, fmt.Errorf("文章评论类型非法")
	}

	// TODO 提取文章摘要

	// 提交数据库
	_ = dao.Begin()
	if post.Pid == 0 {
		post.Pid, v = dao.Post().Add(post)
		if !v {
			_ = dao.Rollback()
			return 0, fmt.Errorf("新建文章失败")
		}
		v = dao.Post().Meta().Add(post.Pid, "post_view_count", "0")
		if !v {
			_ = dao.Rollback()
			return 0, fmt.Errorf("新建文章附加字段失败")
		}
	} else {
		v = dao.Post().Modify(post)
		if !v {
			_ = dao.Rollback()
			return 0, fmt.Errorf("修改文章失败")
		}
	}

	// 处理分类、标签关系
	if post.Pid > 0 {
		v = dao.Rela().RemoveByPostID(post.Pid)
		if !v {
			_ = dao.Rollback()
			return 0, fmt.Errorf("删除旧关系失败")
		}
	}
	tids := make([]uint32, 0, 1+len(req.TagsID))
	tids = append(tids, req.CateID)
	tids = append(tids, req.TagsID...)
	v = dao.Rela().AddMore(post.Pid, tids...)
	if !v {
		_ = dao.Rollback()
		return 0, fmt.Errorf("建立新关系失败")
	}

	// 提交更改
	err := dao.Commit()
	if err != nil {
		return 0, fmt.Errorf("提交事务失败")
	}
	return post.Pid, nil
}

// 修改文章状态
func ModifyPostStatus(req *request.PostStatusModify) error {
	dao := dao.NewDao()
	var v bool
	// 校验文章是否存在
	_, v = dao.Post().GetByID(req.ID)
	if !v {
		return fmt.Errorf("文章不存在")
	}

	// 校验文章类型是否合法
	t := model.PostStatusType(req.Status)
	if t != model.POST_STATUS_PUBLISH && t != model.POST_STATUS_STICKY && t != model.POST_STATUS_DELETED {
		return fmt.Errorf("文章类型非法")
	}

	// 提交数据库
	_ = dao.Begin()
	v = dao.Post().ModifyStatus(req.ID, t)
	if !v {
		_ = dao.Rollback()
		return fmt.Errorf("文章类型修改失败")
	}

	err := dao.Commit()
	if err != nil {
		return fmt.Errorf("事务提交失败")
	}

	return nil
}

// 获取单个文章，返回文章模型（不包含分类、标签）
func GetPost(id uint32) (*model.Post, error) {
	dao := dao.NewDao()
	// 获取文章
	post, v := dao.Post().GetByID(id)
	if !v {
		return nil, fmt.Errorf("文章不存在")
	}
	return post, nil
}

// 获取单个文章的分类、标签，返回分类ID、标签ID列表
func GetPostCateTags(id uint32) (uint32, []uint32, error) {
	dao := dao.NewDao()
	// 获取分类项
	cate, v := dao.Term().ListByPostID(id, model.TERM_TYPE_CATEGORY)
	if !v {
		return 0, nil, fmt.Errorf("获取分类失败")
	}

	// 获取标签项
	tags, _ := dao.Term().ListByPostID(id, model.TERM_TYPE_TAG)

	// 整理
	cateID := cate[0].Tid
	tagsID := make([]uint32, len(tags))
	for i, val := range tags {
		tagsID[i] = val.Tid
	}
	return cateID, tagsID, nil
}
