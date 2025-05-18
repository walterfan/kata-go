package models

import "gorm.io/gorm"

type Prompt struct {
    ID           uint           `json:"id" gorm:"primaryKey"`
    Name         string         `json:"name" gorm:"index"`
    Description  string         `json:"desc"`
    SystemPrompt string         `json:"systemPrompt"`
    UserPrompt   string         `json:"userPrompt"`
    Tags         string         `json:"tags"` // 用逗号分隔的 tag 字符串
    CreatedAt    int64          `json:"createdAt" gorm:"autoCreateTime"`
    UpdatedAt    int64          `json:"updatedAt" gorm:"autoUpdateTime"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
