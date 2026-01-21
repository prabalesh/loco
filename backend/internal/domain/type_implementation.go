package domain

import "time"

type TypeImplementation struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	CustomTypeID     int       `json:"custom_type_id" gorm:"not null;index"`
	LanguageID       int       `json:"language_id" gorm:"not null;index"`
	ClassDefinition  string    `json:"class_definition" gorm:"type:text;not null"`  // TreeNode class code
	SerializerCode   string    `json:"serializer_code" gorm:"type:text;not null"`   // Convert TreeNode -> array
	DeserializerCode string    `json:"deserializer_code" gorm:"type:text;not null"` // Convert array -> TreeNode
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	CustomType CustomType `json:"custom_type,omitempty" gorm:"foreignKey:CustomTypeID;constraint:OnDelete:CASCADE"`
	Language   Language   `json:"language,omitempty" gorm:"foreignKey:LanguageID;constraint:OnDelete:CASCADE"`
}
