package request

import (
	"encoding/json"
	"time"
)

// PageRequest 分页请求
type PageRequest struct {
	Page     int `json:"page" form:"page" binding:"min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"min=1,max=100"`
}

// OptionalTime 用于区分 JSON 中「字段未出现」与「显式传 null」两种情况。
// 在部分更新（PATCH 语义）接口里，前者应保留原值，后者应清空，
// 普通指针的 nil 无法区分二者，故用一个显式的 Set 标记。
type OptionalTime struct {
	Value *time.Time
	Set   bool
}

// UnmarshalJSON 记录字段是否出现，并解析时间值（null 视为清空）
func (o *OptionalTime) UnmarshalJSON(data []byte) error {
	o.Set = true
	if string(data) == "null" {
		o.Value = nil
		return nil
	}
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	o.Value = &t
	return nil
}

// GetOffset 获取偏移量
func (p *PageRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.GetPageSize()
}

// GetPageSize 获取每页数量
func (p *PageRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}
