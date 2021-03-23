package service

import (
	"fmt"
	"goblog/repository"
	"goblog/utils"
	"goblog/utils/errors"
	"os"
	"path"
	"strings"
	"time"
)

func (s *Service) Upload(filename string) (imgPath string, err error) {
	now := time.Now()
	// when the dir used to save img not exist, make it
	dir := fmt.Sprintf("./static/upload/%s/", now.Format("2006/01"))
	if !utils.Exist(dir) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			err = errors.WrapErrorWithStack(err)
			return
		}
	}

	// verify if it's acceptalbe file
	ext := path.Ext(filename)
	if ext != ".jpg" && ext != ".png" && ext != ".gif" && ext != ".bmp" && ext != ".jpeg" {
		err = errors.NewErrorWithStack(fmt.Sprintf("don't accept file with the ext '%s'", ext))
		return
	}

	filenameWithoutExt := strings.TrimSuffix(filename, ext)

	// find the filename not used
	imgPath = fmt.Sprintf("%s%s", dir, filename)
	for t, i := true, 1; t; {
		if utils.Exist(imgPath) {
			i++
			imgPath = fmt.Sprintf("%s%s(%d)%s", dir, filenameWithoutExt, i, ext)
		} else {
			t = false
		}
	}

	s.repository.Begin()
	s.repository.InsertAttachment(&repository.Attachment{
		Created:      uint32(now.Unix()),
		RelativePath: strings.Replace(imgPath, "./static/upload/", "", 1),
	})
	s.repository.Commit()

	err = s.repository.GetError()
	return
}
