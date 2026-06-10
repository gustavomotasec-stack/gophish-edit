package models

import (
	"errors"
	"time"
)

// PhoneList is a collection of phone numbers for WhatsApp campaigns.
type PhoneList struct {
	Id           int64         `json:"id" gorm:"column:id; primary_key:yes"`
	UserId       int64         `json:"-" gorm:"column:user_id"`
	Name         string        `json:"name"`
	Numbers      []PhoneNumber `json:"numbers" gorm:"foreignkey:PhoneListId"`
	ModifiedDate time.Time     `json:"modified_date"`
}

// PhoneNumber is a single entry in a PhoneList.
type PhoneNumber struct {
	Id          int64  `json:"id" gorm:"column:id; primary_key:yes"`
	PhoneListId int64  `json:"-" gorm:"column:phone_list_id"`
	Number      string `json:"number"`
}

func (PhoneList) TableName() string   { return "phone_lists" }
func (PhoneNumber) TableName() string { return "phone_numbers" }

var ErrPhoneListNameNotSpecified = errors.New("Phone list name not specified")
var ErrPhoneListEmpty = errors.New("Phone list must have at least one number")

func (p *PhoneList) Validate() error {
	if p.Name == "" {
		return ErrPhoneListNameNotSpecified
	}
	if len(p.Numbers) == 0 {
		return ErrPhoneListEmpty
	}
	return nil
}

func GetPhoneLists(uid int64) ([]PhoneList, error) {
	var lists []PhoneList
	err := db.Where("user_id = ?", uid).Find(&lists).Error
	if err != nil {
		return lists, err
	}
	for i := range lists {
		db.Where("phone_list_id = ?", lists[i].Id).Find(&lists[i].Numbers)
	}
	return lists, nil
}

func GetPhoneList(id int64, uid int64) (PhoneList, error) {
	var p PhoneList
	err := db.Where("user_id = ? AND id = ?", uid, id).First(&p).Error
	if err != nil {
		return p, err
	}
	db.Where("phone_list_id = ?", p.Id).Find(&p.Numbers)
	return p, nil
}

func PostPhoneList(p *PhoneList) error {
	if err := p.Validate(); err != nil {
		return err
	}
	p.ModifiedDate = time.Now().UTC()
	if err := db.Save(p).Error; err != nil {
		return err
	}
	for i := range p.Numbers {
		p.Numbers[i].PhoneListId = p.Id
		db.Save(&p.Numbers[i])
	}
	return nil
}

func PutPhoneList(p *PhoneList) error {
	if err := p.Validate(); err != nil {
		return err
	}
	p.ModifiedDate = time.Now().UTC()
	db.Where("phone_list_id = ?", p.Id).Delete(&PhoneNumber{})
	if err := db.Save(p).Error; err != nil {
		return err
	}
	for i := range p.Numbers {
		p.Numbers[i].PhoneListId = p.Id
		p.Numbers[i].Id = 0
		db.Save(&p.Numbers[i])
	}
	return nil
}

func DeletePhoneList(id int64, uid int64) error {
	db.Where("phone_list_id = ?", id).Delete(&PhoneNumber{})
	return db.Where("user_id = ? AND id = ?", uid, id).Delete(&PhoneList{}).Error
}
