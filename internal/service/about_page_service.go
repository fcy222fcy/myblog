package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/pkg/logger"
	"encoding/json"
	"sort"

	"go.uber.org/zap"
)

type aboutPageService struct {
	aboutPageRepo repository.AboutPageRepository
}

// NewAboutPageService 创建关于页面服务
func NewAboutPageService(aboutPageRepo repository.AboutPageRepository) AboutPageService {
	return &aboutPageService{aboutPageRepo: aboutPageRepo}
}

// sortSiteHistory 按日期排序建站历程（日期升序，最早的在前）
func sortSiteHistory(siteHistoryJSON string) string {
	var items []entity.SiteHistoryItem
	if err := json.Unmarshal([]byte(siteHistoryJSON), &items); err != nil {
		return siteHistoryJSON
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Date < items[j].Date
	})
	data, err := json.Marshal(items)
	if err != nil {
		return siteHistoryJSON
	}
	return string(data)
}

// GetAboutPage 获取关于页面
func (s *aboutPageService) GetAboutPage() (*request.AboutPageResponse, error) {
	page, err := s.aboutPageRepo.Get()
	if err != nil {
		logger.Warn("获取关于页面失败", zap.Error(err))
		return nil, err
	}

	if page == nil {
		// 返回默认数据
		return &request.AboutPageResponse{
			Title:       "关于我",
			Subtitle:    "Go 开发者 / 独立游戏开发者",
			Bio:         "来自中国的程序员，擅长游戏和后端开发，喜欢玩游戏和骑自行车。",
			Skills:      `["Golang","Erlang","Unity","Docker"]`,
			AboutMe:     `[{"label":"职业","value":"Go 开发者","icon":"💼"},{"label":"爱好","value":"游戏 / 骑行","icon":"🎮"},{"label":"技术栈","value":"Go / Erlang / Unity","icon":"🛠"},{"label":"邮箱","value":"3181088318@qq.com","icon":"✉️"}]`,
			AboutSite:   `[{"label":"框架","value":"Go + Gin","icon":"⚙️"},{"label":"前端","value":"Vue 3","icon":"🎨"},{"label":"部署","value":"Docker","icon":"☁️"}]`,
			Projects:    `[{"name":"GitHub","description":"我的 GitHub 主页，包含各种编程项目。","url":"https://github.com/fcy222fcy","icon":"⌘"},{"name":"DesktopSnap","description":"Windows 桌面图标保存和恢复工具。","url":"https://desktopsnap.fcyan.xyz/","icon":"🖥"}]`,
			ContactInfo: `[{"label":"GitHub","value":"fcy222fcy","icon":"🐙","url":"https://github.com/fcy222fcy"}]`,
			SiteHistory: `[]`,
		}, nil
	}

	return &request.AboutPageResponse{
		ID:          page.ID,
		Title:       page.Title,
		Subtitle:    page.Subtitle,
		Bio:         page.Bio,
		Skills:      page.Skills,
		AboutMe:     page.AboutMe,
		AboutSite:   page.AboutSite,
		Projects:    page.Projects,
		ContactInfo: page.ContactInfo,
		SiteHistory: sortSiteHistory(page.SiteHistory),
	}, nil
}

// UpdateAboutPage 更新关于页面
func (s *aboutPageService) UpdateAboutPage(req *request.UpdateAboutPageRequest) error {
	page, err := s.aboutPageRepo.Get()
	if err != nil {
		return err
	}

	if page == nil {
		// 创建新记录
		page = &entity.AboutPage{}
	}

	page.Title = req.Title
	page.Subtitle = req.Subtitle
	page.Bio = req.Bio
	page.Skills = req.Skills
	page.AboutMe = req.AboutMe
	page.AboutSite = req.AboutSite
	page.Projects = req.Projects
	page.ContactInfo = req.ContactInfo
	page.SiteHistory = sortSiteHistory(req.SiteHistory)

	return s.aboutPageRepo.Save(page)
}
