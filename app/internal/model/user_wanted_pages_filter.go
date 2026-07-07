package model

import "github.com/lib/pq"

type UserWantedPagesFilter struct {
	FilterData UserWantedPageCriteria `json:"filterData" gorm:"embedded"`
	User       User                    `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Model
	UserID uint `json:"userId" gorm:"not null;column:user_id;"`
}

type UserWantedPageCriteria struct {
	Name          string         `json:"name" gorm:"column:name;not null;default:''"`
	SiteNames     pq.StringArray `json:"siteNames,omitempty" gorm:"column:site_names;type:text[]"`
	Keywords      pq.StringArray `json:"keywords,omitempty" gorm:"column:keywords;type:text[]"`
	TitleContains string         `json:"titleContains,omitempty" gorm:"column:title_contains"`
}
