package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	ctx "github.com/gophish/gophish/context"
	"github.com/gophish/gophish/models"
	"github.com/gorilla/mux"
)

// WhatsAppTemplates handles GET/POST /api/whatsapp_templates/
func (as *Server) WhatsAppTemplates(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Get(r, "user_id").(int64)
	switch r.Method {
	case http.MethodGet:
		ts, err := models.GetWhatsAppTemplates(uid)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, ts, http.StatusOK)
	case http.MethodPost:
		t := models.WhatsAppTemplate{}
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		t.UserId = uid
		if err := models.PostWhatsAppTemplate(&t); err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		JSONResponse(w, t, http.StatusCreated)
	}
}

// WhatsAppTemplate handles GET/PUT/DELETE /api/whatsapp_templates/{id}
func (as *Server) WhatsAppTemplate(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Get(r, "user_id").(int64)
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	switch r.Method {
	case http.MethodGet:
		t, err := models.GetWhatsAppTemplate(id, uid)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Template not found"}, http.StatusNotFound)
			return
		}
		JSONResponse(w, t, http.StatusOK)
	case http.MethodPut:
		t := models.WhatsAppTemplate{}
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		t.Id = id
		t.UserId = uid
		if err := models.PutWhatsAppTemplate(&t); err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		JSONResponse(w, t, http.StatusOK)
	case http.MethodDelete:
		if err := models.DeleteWhatsAppTemplate(id, uid); err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "Template deleted"}, http.StatusOK)
	}
}
