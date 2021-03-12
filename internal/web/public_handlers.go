package web

import (
	"fmt"
	"log"
	"net/http"
)

func (inst *instance) getIndex(w http.ResponseWriter, r *http.Request) {
	contactSessionData, err := inst.generateContactSessionDataForTemplate()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = inst.indexTemplate.Execute(w, contactSessionData)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (inst *instance) contact(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		contactSessionData, err := inst.generateContactSessionDataForTemplate()
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = inst.contactTemplate.Execute(w, contactSessionData)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case "POST":
		validContact, formErrors, formData := inst.validContactForm(w, r)
		if !validContact {
			contactSessionData, err := inst.generateContactSessionDataForTemplate()
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			contactSessionData["formErrors"] = formErrors
			contactSessionData["contactName"] = formData.contactName
			contactSessionData["contactEmail"] = formData.contactEmail
			contactSessionData["contactMessage"] = formData.contactMessage
			err = inst.contactTemplate.Execute(w, contactSessionData)
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			return
		}

		subject := fmt.Sprintf("[clinton.dev Contact] Message from %s", formData.contactEmail)
		err := inst.email.SendEmail(inst.contactEmailFrom, formData.contactEmail,
			inst.contactEmailTo, subject, formData.contactMessage)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		contactSessionData, err := inst.generateContactSessionDataForTemplate()
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		contactSessionData["contactName"] = formData.contactName
		contactSessionData["formSuccess"] = true
		err = inst.contactTemplate.Execute(w, contactSessionData)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
