package models

import (
	"crypto/rand"
	"math/big"
	"time"
)

const shortLinkChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const shortLinkLen = 7

// ShortLink maps a short code to a full phishing URL.
type ShortLink struct {
	Id         int64     `json:"id" gorm:"primary_key"`
	Code       string    `json:"code"`
	Original   string    `json:"original"`
	CampaignId int64     `json:"campaign_id"`
	RId        string    `json:"rid"`
	CreatedAt  time.Time `json:"created_at"`
}

func generateShortCode() (string, error) {
	b := make([]byte, shortLinkLen)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(shortLinkChars))))
		if err != nil {
			return "", err
		}
		b[i] = shortLinkChars[n.Int64()]
	}
	return string(b), nil
}

// CreateShortLink generates a unique code and stores the mapping.
func CreateShortLink(original string, campaignId int64, rid string) (ShortLink, error) {
	// Check if already exists for this rid
	var existing ShortLink
	if err := db.Where("rid = ? AND campaign_id = ?", rid, campaignId).First(&existing).Error; err == nil {
		return existing, nil
	}

	var sl ShortLink
	for {
		code, err := generateShortCode()
		if err != nil {
			return sl, err
		}
		var count int
		db.Model(&ShortLink{}).Where("code = ?", code).Count(&count)
		if count == 0 {
			sl = ShortLink{
				Code:       code,
				Original:   original,
				CampaignId: campaignId,
				RId:        rid,
				CreatedAt:  time.Now().UTC(),
			}
			err = db.Create(&sl).Error
			return sl, err
		}
	}
}

// GetShortLink returns the ShortLink for a given code.
func GetShortLink(code string) (ShortLink, error) {
	var sl ShortLink
	err := db.Where("code = ?", code).First(&sl).Error
	return sl, err
}
