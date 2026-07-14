package model

type SubscriptionLevelProduct struct {
	SubscriptionLevel   SubscriptionLevel `json:"subscriptionLevel" gorm:"foreignKey:SubscriptionLevelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Product             Product           `json:"product" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SubscriptionLevelID int               `json:"subscriptionLevelId" gorm:"primaryKey;not null"`
	ProductID           int               `json:"productId" gorm:"primaryKey;not null"`
}
