package dto

type PageDTO struct {
	SiteName   string `json:"site_name"`
	URL        string `json:"url"`
	Content    string `json:"content"`
	Title      string `json:"title"`
	TextField1 string `json:"text_field_1"`
	TextField2 string `json:"text_field_2"`
	TextField3 string `json:"text_field_3"`
	ExternalID string `json:"external_id"`
}
