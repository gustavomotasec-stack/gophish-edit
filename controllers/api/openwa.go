package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	ctx "github.com/gophish/gophish/context"
	log "github.com/gophish/gophish/logger"
	"github.com/gophish/gophish/models"
)

// OpenWAConfig handles GET/POST /api/openwa/config
func (as *Server) OpenWAConfig(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Get(r, "user_id").(int64)
	switch r.Method {
	case http.MethodGet:
		cfg, err := models.GetOpenWAConfig(uid)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, cfg, http.StatusOK)
	case http.MethodPost:
		cfg := models.OpenWAConfig{}
		if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		cfg.UserId = uid
		if cfg.MinDelay <= 0 {
			cfg.MinDelay = 3
		}
		if cfg.MaxDelay <= cfg.MinDelay {
			cfg.MaxDelay = cfg.MinDelay + 5
		}
		if err := models.SaveOpenWAConfig(&cfg); err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, cfg, http.StatusOK)
	}
}

// WADispatchRequest is the payload for POST /api/openwa/dispatch
type WADispatchRequest struct {
	TemplateId  int64     `json:"template_id"`
	PhoneListId int64     `json:"phone_list_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	CampaignId  int64     `json:"campaign_id"`
}

// OpenWADispatch handles POST /api/openwa/dispatch
// It spawns a goroutine that sends messages with random delays.
func (as *Server) OpenWADispatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	uid := ctx.Get(r, "user_id").(int64)

	var req WADispatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
		return
	}

	tmpl, err := models.GetWhatsAppTemplate(req.TemplateId, uid)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Template não encontrado"}, http.StatusNotFound)
		return
	}

	phoneList, err := models.GetPhoneList(req.PhoneListId, uid)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Lista de números não encontrada"}, http.StatusNotFound)
		return
	}

	cfg, err := models.GetOpenWAConfig(uid)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	// Calculate phish base URL from campaign
	baseURL := ""
	if req.CampaignId > 0 {
		campaign, cerr := models.GetCampaign(req.CampaignId, uid)
		if cerr == nil {
			baseURL = campaign.URL
		}
	}

	totalNumbers := len(phoneList.Numbers)
	JSONResponse(w, models.Response{
		Success: true,
		Message: fmt.Sprintf("Disparo iniciado para %d números", totalNumbers),
	}, http.StatusOK)

	// Dispatch asynchronously
	go func() {
		delay := !req.ScheduledAt.IsZero() && req.ScheduledAt.After(time.Now())
		if delay {
			waitDur := time.Until(req.ScheduledAt)
			log.Infof("WhatsApp dispatch scheduled in %v", waitDur)
			time.Sleep(waitDur)
		}

		for i, num := range phoneList.Numbers {
			// Generate short link for this number
			rid := fmt.Sprintf("wa%d_%d", req.CampaignId, num.Id)
			origURL := fmt.Sprintf("%s?rid=%s", baseURL, rid)
			shortURL := origURL

			if baseURL != "" {
				sl, serr := models.CreateShortLink(origURL, req.CampaignId, rid)
				if serr == nil {
					// Extract base from origURL
					shortURL = extractBase(origURL) + "/r/" + sl.Code
				}
			}

			message := models.RenderWAMessage(tmpl.Body, "", "", shortURL)

			if err := models.OpenWASendMessage(cfg, num.Number, message); err != nil {
				log.Errorf("Failed to send WhatsApp to %s: %v", num.Number, err)
			} else {
				log.Infof("WhatsApp sent to %s (%d/%d)", num.Number, i+1, totalNumbers)
			}

			// Random delay between messages (skip after last)
			if i < totalNumbers-1 {
				models.RandomDelay(cfg.MinDelay, cfg.MaxDelay)
			}
		}
		log.Infof("WhatsApp dispatch completed: %d messages", totalNumbers)
	}()
}
