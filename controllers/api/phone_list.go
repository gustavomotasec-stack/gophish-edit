package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	ctx "github.com/gophish/gophish/context"
	"github.com/gophish/gophish/models"
	"github.com/gorilla/mux"
)

// PhoneLists handles GET/POST /api/phone_lists/
func (as *Server) PhoneLists(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Get(r, "user_id").(int64)
	switch r.Method {
	case http.MethodGet:
		lists, err := models.GetPhoneLists(uid)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, lists, http.StatusOK)
	case http.MethodPost:
		p := models.PhoneList{}
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		p.UserId = uid
		if err := models.PostPhoneList(&p); err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		JSONResponse(w, p, http.StatusCreated)
	}
}

// PhoneListById handles GET/PUT/DELETE /api/phone_lists/{id}
func (as *Server) PhoneListById(w http.ResponseWriter, r *http.Request) {
	uid := ctx.Get(r, "user_id").(int64)
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	switch r.Method {
	case http.MethodGet:
		p, err := models.GetPhoneList(id, uid)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: "Phone list not found"}, http.StatusNotFound)
			return
		}
		JSONResponse(w, p, http.StatusOK)
	case http.MethodPut:
		p := models.PhoneList{}
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		p.Id = id
		p.UserId = uid
		if err := models.PutPhoneList(&p); err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusBadRequest)
			return
		}
		JSONResponse(w, p, http.StatusOK)
	case http.MethodDelete:
		if err := models.DeletePhoneList(id, uid); err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "Phone list deleted"}, http.StatusOK)
	}
}
