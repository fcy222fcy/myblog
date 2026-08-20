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

// categoryService 分类服务实现
type categoryService struct {
	categoryRepo repository.CategoryRepository
	redisClient  *redis.Client
}

// NewCategoryService 创建分类服务
func NewCategoryService(categoryRepo repository.CategoryRepository, redisClient *redis.Client) CategoryService {
	return &categoryService{categoryRepo: categoryRepo, redisClient: redisClient}
}

// GetCategoryList 获取分类列表
func (s *categoryService) GetCategoryList() ([]response.CategoryResponse, error) {
	if s.redisClient != nil {
		var cachedList []response.CategoryResponse
		if err := s.redisClient.GetJSON(context.Background(), "category:list", &cachedList); err == nil {
			return cachedList, nil
		}
	}

	categories, err := s.categoryRepo.List()
	if err != nil {
		return nil, fmt.Errorf("获取分类列表失败, %w", err)
	}

	var result []response.CategoryResponse
	for _, cat := range categories {
		count, _ := s.categoryRepo.GetCategoryArticleCount(cat.ID)
		result = append(result, response.CategoryResponse{
			ID:           cat.ID,
			Name:         cat.Name,
			Slug:         cat.Slug,
			Description:  cat.Description,
			Icon:         cat.Icon,
			SortOrder:    cat.SortOrder,
			ArticleCount: count,
		})
	}

	if s.redisClient != nil {
		go s.redisClient.SetJSON(context.Background(), "category:list", result, 30*time.Minute)
	}

	return result, nil
}

// GetCategoryByID 根据ID获取分类
func (s *categoryService) GetCategoryByID(id uint) (*response.CategoryResponse, error) {
	cat, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取分类失败, %w", err)
	}
	if cat == nil {
		return nil, bizerrors.New(bizerrors.CodeCategoryNotFound, bizerrors.GetMessage(bizerrors.CodeCategoryNotFound))
	}

	count, _ := s.categoryRepo.GetCategoryArticleCount(cat.ID)
	return &response.CategoryResponse{
		ID:           cat.ID,
		Name:         cat.Name,
		Slug:         cat.Slug,
		Description:  cat.Description,
		Icon:         cat.Icon,
		SortOrder:    cat.SortOrder,
		ArticleCount: count,
	}, nil
}

// CreateCategory 创建分类
func (s *categoryService) CreateCategory(req *request.CreateCategoryRequest) (uint, error) {
	existing, err := s.categoryRepo.FindByName(req.Name)
	if err != nil {
		return 0, fmt.Errorf("查询分类失败, %w", err)
	}
	if existing != nil {
		return 0, bizerrors.New(bizerrors.CodeCategoryNameExists, bizerrors.GetMessage(bizerrors.CodeCategoryNameExists))
	}

	// 如果 slug 为空，自动生成
	categorySlug := req.Slug
	if categorySlug == "" {
		categorySlug = slug.Generate(req.Name)
		// 确保 slug 唯一
		existingSlug, _ := s.categoryRepo.FindBySlug(categorySlug)
		if existingSlug != nil {
			categorySlug = fmt.Sprintf("%s-%d", categorySlug, time.Now().Unix())
		}
	}

	category := &entity.Category{
		Name:        req.Name,
		Slug:        categorySlug,
		Description: req.Description,
		Icon:        req.Icon,
		SortOrder:   req.SortOrder,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return 0, fmt.Errorf("创建分类失败, %w", err)
	}

	logger.Infof("创建分类成功, id=%d, name=%s", category.ID, category.Name)

	s.invalidateCategoryCache()

	return category.ID, nil
}

// UpdateCategory 更新分类
func (s *categoryService) UpdateCategory(id uint, req *request.UpdateCategoryRequest) error {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("查询分类失败, %w", err)
	}
	if category == nil {
		return bizerrors.New(bizerrors.CodeCategoryNotFound, bizerrors.GetMessage(bizerrors.CodeCategoryNotFound))
	}

	if req.Name != "" {
		existing, err := s.categoryRepo.FindByName(req.Name)
		if err != nil {
			return fmt.Errorf("查询分类名称失败, %w", err)
		}
		if existing != nil && existing.ID != id {
			return bizerrors.New(bizerrors.CodeCategoryNameExists, bizerrors.GetMessage(bizerrors.CodeCategoryNameExists))
		}
		category.Name = req.Name
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}
	if req.Description != nil {
		category.Description = *req.Description
	}
	if req.Icon != nil {
		category.Icon = *req.Icon
	}
	if req.SortOrder != 0 {
		category.SortOrder = req.SortOrder
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return fmt.Errorf("更新分类失败, %w", err)
	}

	logger.Infof("更新分类成功, id=%d", id)

	s.invalidateCategoryCache()

	return nil
}

// DeleteCategory 删除分类
func (s *categoryService) DeleteCategory(id uint) error {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("查询分类失败, %w", err)
	}
	if category == nil {
		return bizerrors.New(bizerrors.CodeCategoryNotFound, bizerrors.GetMessage(bizerrors.CodeCategoryNotFound))
	}

	count, err := s.categoryRepo.GetCategoryArticleCount(id)
	if err != nil {
		return fmt.Errorf("获取分类文章数失败, %w", err)
	}
	if count > 0 {
		return bizerrors.New(bizerrors.CodeCategoryHasArticles, bizerrors.GetMessage(bizerrors.CodeCategoryHasArticles))
	}

	if err := s.categoryRepo.Delete(id); err != nil {
		return fmt.Errorf("删除分类失败, %w", err)
	}

	logger.Infof("删除分类成功, id=%d", id)

	s.invalidateCategoryCache()

	return nil
}

func (s *categoryService) invalidateCategoryCache() {
	if s.redisClient == nil {
		return
	}
	ctx := context.Background()
	keys := []string{
		"category:list",
	}
	go s.redisClient.Del(ctx, keys...)
}
