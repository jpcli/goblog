package service

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"goblog/cache"
)

func (s *Service) GetWebsiteBasicOption() (map[string]string, error) {
	res := make(map[string]string)
	s.getOption("websiteURL", res)
	s.getOption("websiteName", res)
	s.getOption("eachPageLimit", res)
	s.getOption("pageNavLimit", res)
	return res, s.repository.GetError()
}

func (s *Service) getOption(key string, dest map[string]string) {
	dest[key] = s.repository.GetOption(key)
}

func (s *Service) SetWebsiteBasicOption(websiteURL, websiteName string, eachPageLimit, pageNavLimit uint32) error {
	s.repository.Begin()
	s.repository.SetOption("websiteURL", websiteURL)
	s.repository.SetOption("websiteName", websiteName)
	s.repository.SetOption("eachPageLimit", cast.ToString(eachPageLimit))
	s.repository.SetOption("pageNavLimit", cast.ToString(pageNavLimit))
	s.repository.Commit()
	if err := s.repository.GetError(); err != nil {
		return err
	} else {
		cache.NewCache().ClearOption()
		return nil
	}

}

func (s *Service) setTKD(prefix, t, k, d string) {
	b, _ := json.Marshal(map[string]string{
		"title":       t,
		"keywords":    k,
		"description": d,
	})
	d = string(b)
	s.repository.SetOption(fmt.Sprintf("%sTKD", prefix), d)
}
