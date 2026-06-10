package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// OpenWAConfig stores the connection settings for OpenWA.
type OpenWAConfig struct {
	Id       int64  `json:"id" gorm:"column:id; primary_key:yes"`
	UserId   int64  `json:"-" gorm:"column:user_id"`
	APIURL   string `json:"api_url"`
	APIKey   string `json:"api_key"`
	MinDelay int    `json:"min_delay"`
	MaxDelay int    `json:"max_delay"`
}

func (OpenWAConfig) TableName() string { return "openwa_configs" }

func GetOpenWAConfig(uid int64) (OpenWAConfig, error) {
	var c OpenWAConfig
	err := db.Where("user_id = ?", uid).First(&c).Error
	if err != nil {
		// Return defaults if not configured yet
		return OpenWAConfig{UserId: uid, APIURL: "http://localhost:3000", MinDelay: 3, MaxDelay: 8}, nil
	}
	return c, nil
}

func SaveOpenWAConfig(c *OpenWAConfig) error {
	var existing OpenWAConfig
	if db.Where("user_id = ?", c.UserId).First(&existing).Error == nil {
		c.Id = existing.Id
	}
	return db.Save(c).Error
}

// OpenWASendMessage sends a single WhatsApp message via the OpenWA API.
func OpenWASendMessage(cfg OpenWAConfig, phone, message string) error {
	// Normalize phone: ensure it ends with @c.us for OpenWA
	if !strings.Contains(phone, "@") {
		phone = phone + "@c.us"
	}

	payload := map[string]string{
		"chatId":  phone,
		"content": message,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", cfg.APIURL+"/api/sendText", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.APIKey != "" {
		req.Header.Set("Authorization", cfg.APIKey)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("OpenWA returned status %d", resp.StatusCode)
	}
	return nil
}

// RenderWAMessage substitutes template variables in a WhatsApp message body.
func RenderWAMessage(body string, firstName, lastName, shortURL string) string {
	r := strings.NewReplacer(
		"{{.FirstName}}", firstName,
		"{{.LastName}}", lastName,
		"{{.URL}}", shortURL,
	)
	return r.Replace(body)
}

// RandomDelay sleeps for a random duration between min and max seconds.
func RandomDelay(min, max int) {
	if min <= 0 {
		min = 1
	}
	if max <= min {
		max = min + 1
	}
	n := min + rand.Intn(max-min)
	time.Sleep(time.Duration(n) * time.Second)
}
