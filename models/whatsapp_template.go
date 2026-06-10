package models

import (
	"errors"
	"time"
)

// WhatsAppTemplate holds a reusable WhatsApp message template.
type WhatsAppTemplate struct {
	Id           int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserId       int64     `json:"-" gorm:"column:user_id"`
	Name         string    `json:"name"`
	Body         string    `json:"body"`
	ModifiedDate time.Time `json:"modified_date"`
}

// TableName overrides the default table name.
func (WhatsAppTemplate) TableName() string { return "whatsapp_templates" }

var ErrWATemplateNameNotSpecified = errors.New("Template name not specified")
var ErrWATemplateBodyNotSpecified = errors.New("Template body not specified")

func (t *WhatsAppTemplate) Validate() error {
	if t.Name == "" {
		return ErrWATemplateNameNotSpecified
	}
	if t.Body == "" {
		return ErrWATemplateBodyNotSpecified
	}
	return nil
}

func GetWhatsAppTemplates(uid int64) ([]WhatsAppTemplate, error) {
	var ts []WhatsAppTemplate
	err := db.Where("user_id = ?", uid).Find(&ts).Error
	return ts, err
}

func GetWhatsAppTemplate(id int64, uid int64) (WhatsAppTemplate, error) {
	var t WhatsAppTemplate
	err := db.Where("user_id = ? AND id = ?", uid, id).First(&t).Error
	return t, err
}

func PostWhatsAppTemplate(t *WhatsAppTemplate) error {
	if err := t.Validate(); err != nil {
		return err
	}
	t.ModifiedDate = time.Now().UTC()
	return db.Save(t).Error
}

func PutWhatsAppTemplate(t *WhatsAppTemplate) error {
	if err := t.Validate(); err != nil {
		return err
	}
	t.ModifiedDate = time.Now().UTC()
	return db.Save(t).Error
}

func DeleteWhatsAppTemplate(id int64, uid int64) error {
	return db.Where("user_id = ? AND id = ?", uid, id).Delete(&WhatsAppTemplate{}).Error
}
