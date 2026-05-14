package model

type Page struct {
	Model
	SiteName   string `json:"siteName" gorm:"not null;column:site_name;index:idx_page_site_external"`
	URL        string `json:"url" gorm:"not null;column:url;uniqueIndex:idx_page_site_external"`
	Content    string `json:"content" gorm:"type:text;column:content"`
	Title      string `json:"title" gorm:"not null;column:title"`
	TextField1 string `json:"textField1" gorm:"column:text_field_1"`
	TextField2 string `json:"textField2" gorm:"column:text_field_2"`
	TextField3 string `json:"textField3" gorm:"column:text_field_3"`
	HashKey    string `json:"hashKey" gorm:"not null;column:hash_key;uniqueIndex:idx_page_site_hash"`
	ExternalID string `json:"externalId" gorm:"not null;column:external_id;index:idx_page_site_external"`
	Notified   bool   `json:"notified" gorm:"not null;column:notified;default:false"`
}
