package service

import (
	"fmt"
	"goblog/dao"
	"goblog/model"
	"os"
	"path"
	"strings"
	"time"
)

// 保存附件，返回该文件保存的路径
func AddAttachment(filename string) (string, error) {
	now := time.Now()
	// 创建保存文件夹
	dir := fmt.Sprintf("./static/upload/%s/", now.Format("2006/01"))
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("创建保存文件夹失败")
		}
	}

	// 校验文件后缀是否被允许
	ext := path.Ext(filename)
	if ext != ".jpg" && ext != ".png" && ext != ".gif" && ext != ".bmp" && ext != ".jpeg" {
		return "", fmt.Errorf("该文件不允许被上传")
	}

	// 找到未被使用的文件名
	filenameWithoutExt := strings.TrimSuffix(filename, ext)
	savePath := fmt.Sprintf("%s%s", dir, filename)
	for i := 1; ; {
		if _, err := os.Stat(savePath); err == nil || os.IsExist(err) {
			i++
			savePath = fmt.Sprintf("%s%s(%d)%s", dir, filenameWithoutExt, i, ext)
		} else {
			break
		}
	}

	// 提交数据库
	dao := dao.NewDao()
	_ = dao.Begin()
	_, v := dao.Atta().Add(&model.Attachment{
		Created:      uint32(now.Unix()),
		RelativePath: savePath,
	})
	if !v {
		return "", fmt.Errorf("新建附件失败")
	}
	err := dao.Commit()
	if err != nil {
		return "", fmt.Errorf("提交事务失败")
	}

	return savePath, nil
}
