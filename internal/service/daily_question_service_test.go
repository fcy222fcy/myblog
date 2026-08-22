package service

import (
	"blog/internal/model/entity"
	"blog/internal/repository"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestDailyQuestionReadsAreSideEffectFreeAndAllReturnsFullData(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&entity.DailyQuestion{}))
	question := entity.DailyQuestion{
		Question:  "为什么要把读取和计数拆开？",
		Answer:    "因为 GET 应该保持只读。",
		Date:      "2026-08-22",
		Status:    entity.DailyQuestionStatusPublished,
		ViewCount: 4,
	}
	require.NoError(t, db.Create(&question).Error)
	svc := NewDailyQuestionService(repository.NewDailyQuestionRepository(db))

	all, err := svc.GetAllPublishedQuestions()
	require.NoError(t, err)
	require.Len(t, all, 1)
	require.Equal(t, question.Answer, all[0].Answer)
	require.Equal(t, int64(4), all[0].ViewCount)

	_, err = svc.GetLatestQuestion()
	require.NoError(t, err)
	_, err = svc.GetQuestionByDate(question.Date)
	require.NoError(t, err)

	var stored entity.DailyQuestion
	require.NoError(t, db.First(&stored, question.ID).Error)
	require.Equal(t, int64(4), stored.ViewCount)
}
