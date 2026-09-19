package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	bizerrors "blog/pkg/errors"
	"blog/pkg/logger"
	"blog/pkg/redis"
	"blog/pkg/slug"
	"context"
	"fmt"
	"time"
)

// tagService 标签服务实现
type tagService struct {
	tagRepo     repository.TagRepository
	redisClient *redis.Client
}

// NewTagService 创建标签服务
func NewTagService(tagRepo repository.TagRepository, redisClient *redis.Client) TagService {
	return &tagService{tagRepo: tagRepo, redisClient: redisClient}
}

// GetTagList 获取标签列表
func (s *tagService) GetTagList() ([]response.TagResponse, error) {
	if s.redisClient != nil {
		var cachedList []response.TagResponse
		if err := s.redisClient.GetJSON(context.Background(), "tag:list", &cachedList); err == nil {
			return cachedList, nil
		}
	}

	tags, err := s.tagRepo.List()
	if err != nil {
		return nil, fmt.Errorf("获取标签列表失败, %w", err)
	}

	var result []response.TagResponse
	// 同 GetCategoryList：count 为 0 无法区分「真的没有文章」与「查询失败」，
	// 出错时必须跳过缓存写入，避免错误的 0 被固化进 tag:list 一个缓存周期（30 分钟）。
	countsTrustworthy := true
	for _, tag := range tags {
		count, err := s.tagRepo.GetTagArticleCount(tag.ID)
		if err != nil {
			logger.Warnf("获取标签文章数失败, tagID: %d, error: %v", tag.ID, err)
			countsTrustworthy = false
		}
		result = append(result, response.TagResponse{
			ID:           tag.ID,
			Name:         tag.Name,
			Slug:         tag.Slug,
			ArticleCount: count,
			CreatedAt:    tag.CreatedAt,
		})
	}

	if s.redisClient != nil && countsTrustworthy {
		go s.redisClient.SetJSON(context.Background(), "tag:list", result, 30*time.Minute)
	}

	return result, nil
}

// CreateTag 创建标签
func (s *tagService) CreateTag(req *request.CreateTagRequest) (uint, error) {
	existing, err := s.tagRepo.FindByName(req.Name)
	if err != nil {
		return 0, fmt.Errorf("查询标签失败, %w", err)
	}
	if existing != nil {
		return 0, bizerrors.New(bizerrors.CodeTagNameExists, bizerrors.GetMessage(bizerrors.CodeTagNameExists))
	}

	// 如果 slug 为空，自动生成
	tagSlug := req.Slug
	if tagSlug == "" {
		tagSlug = slug.Generate(req.Name)
		// 确保 slug 唯一
		existingSlug, _ := s.tagRepo.FindBySlug(tagSlug)
		if existingSlug != nil {
			tagSlug = fmt.Sprintf("%s-%d", tagSlug, time.Now().Unix())
		}
	}

	tag := &entity.Tag{
		Name: req.Name,
		Slug: tagSlug,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return 0, fmt.Errorf("创建标签失败, %w", err)
	}

	logger.Infof("创建标签成功, id=%d, name=%s", tag.ID, tag.Name)

	s.invalidateTagCache()

	return tag.ID, nil
}

// UpdateTag 更新标签
func (s *tagService) UpdateTag(id uint, req *request.UpdateTagRequest) error {
	tag, err := s.tagRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("查询标签失败, %w", err)
	}
	if tag == nil {
		return bizerrors.New(bizerrors.CodeTagNotFound, bizerrors.GetMessage(bizerrors.CodeTagNotFound))
	}

	if req.Name != "" {
		existing, err := s.tagRepo.FindByName(req.Name)
		if err != nil {
			return fmt.Errorf("查询标签名称失败, %w", err)
		}
		if existing != nil && existing.ID != id {
			return bizerrors.New(bizerrors.CodeTagNameExists, bizerrors.GetMessage(bizerrors.CodeTagNameExists))
		}
		tag.Name = req.Name
	}
	if req.Slug != nil {
		tag.Slug = *req.Slug
	}

	if err := s.tagRepo.Update(tag); err != nil {
		return fmt.Errorf("更新标签失败, %w", err)
	}

	logger.Infof("更新标签成功, id=%d", id)

	s.invalidateTagCache()

	return nil
}

// DeleteTag 删除标签
func (s *tagService) DeleteTag(id uint) error {
	tag, err := s.tagRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("查询标签失败, %w", err)
	}
	if tag == nil {
		return bizerrors.New(bizerrors.CodeTagNotFound, bizerrors.GetMessage(bizerrors.CodeTagNotFound))
	}

	count, err := s.tagRepo.GetTagArticleCount(id)
	if err != nil {
		return fmt.Errorf("获取标签文章数失败, %w", err)
	}
	if count > 0 {
		return bizerrors.New(bizerrors.CodeTagHasArticles, bizerrors.GetMessage(bizerrors.CodeTagHasArticles))
	}

	if err := s.tagRepo.Delete(id); err != nil {
		return fmt.Errorf("删除标签失败, %w", err)
	}

	logger.Infof("删除标签成功, id=%d", id)

	s.invalidateTagCache()

	return nil
}

func (s *tagService) invalidateTagCache() {
	if s.redisClient == nil {
		return
	}
	ctx := context.Background()
	keys := []string{
		"tag:list",
	}
	go s.redisClient.Del(ctx, keys...)
}
