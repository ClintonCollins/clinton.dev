package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type response struct {
	Success      bool     `json:"success"`
	Errors       []string `json:"errors"`
	NewFormToken string   `json:"new_token"`
	NewFormNum1  int      `json:"new_form_num_1"`
	NewFormNum2  int      `json:"new_form_num_2"`
}

func serveJSONResponse(w http.ResponseWriter, r *http.Request, jsonStruct interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&jsonStruct)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (inst *instance) apiContact(w http.ResponseWriter, r *http.Request) {
	validContact, formErrors, formData := inst.validContactForm(w, r)
	num1, num2, token, err := inst.generateContactSessionData()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if !validContact {
		resp := &response{
			Success:      false,
			Errors:       formErrors,
			NewFormToken: token,
			NewFormNum1:  num1,
			NewFormNum2:  num2,
		}
		serveJSONResponse(w, r, resp)
		return
	}

	subject := fmt.Sprintf("[clinton.dev Contact] Message from %s", formData.contactEmail)
	err = inst.email.SendEmail(inst.contactEmailFrom, formData.contactEmail,
		inst.contactEmailTo, subject, formData.contactMessage)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := &response{
		Success: true,
		Errors:  nil,
		NewFormToken: token,
		NewFormNum1:  num1,
		NewFormNum2:  num2,
	}
	serveJSONResponse(w, r, resp)
}
