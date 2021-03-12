package web

import (
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type contactFormData struct {
	contactName    string
	contactEmail   string
	contactMessage string
}

// parseTemplates just parses all the static templates once when the webserver is initialized instead of on every request.
func (inst *instance) parseTemplates() error {
	t, err := template.ParseFiles(inst.frontendPath+"/templates/layout.html", inst.frontendPath+"/templates/index.html")
	if err != nil {
		return err
	}
	inst.indexTemplate = t
	t, err = template.ParseFiles(inst.frontendPath+"/templates/layout.html",
		inst.frontendPath+"/templates/contact.html")
	if err != nil {
		return err
	}
	inst.contactTemplate = t
	return nil
}

// generateContactSessionData creates a new contact session.
func (inst *instance) generateContactSessionData() (num1, num2 int, token string, err error) {
	uuid4, err := uuid.NewRandom()
	if err != nil {
		return 0, 0, "", err
	}

	numb1 := rand.Intn(10)
	numb2 := rand.Intn(10)

	contactSession := contactFormSession{
		createdAt: time.Now(),
		sumAnswer: numb1 + numb2,
	}

	inst.contactFormSessionsLock.Lock()
	inst.contactFormSessions[uuid4.String()] = contactSession
	inst.contactFormSessionsLock.Unlock()
	return numb1, numb2, uuid4.String(), nil
}

// generateContactSessionDataForTemplate gets a new contact session then bundles it into a map that can be used
// by a template.
func (inst *instance) generateContactSessionDataForTemplate() (map[string]interface{}, error) {
	num1, num2, token, err := inst.generateContactSessionData()
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	data["num1"] = num1
	data["num2"] = num2
	data["contactToken"] = token
	return data, nil
}

// validContactForm validates the contact form either from a template or API call. Returns whether or not it's valid
// along with any errors and the data from the form if valid.
func (inst *instance) validContactForm(w http.ResponseWriter, r *http.Request) (bool, []string, contactFormData) {
	formData := contactFormData{}
	err := r.ParseForm()
	if err != nil {
		return false, []string{"Bad form data."}, formData
	}
	contactName := r.PostFormValue("contactName")
	contactEmail := r.PostFormValue("contactEmail")
	contactMessage := r.PostFormValue("contactMessage")
	contactSumAnswer := r.PostFormValue("contactFormVerificationSum")
	contactSessionToken := r.PostFormValue("contactFormToken")
	honeyPotEmail := r.PostFormValue("contactEmailFake")

	if honeyPotEmail != "" {
		return false, []string{"Bad form data."}, formData
	}

	formData.contactName = contactName
	formData.contactEmail = contactEmail
	formData.contactMessage = contactMessage

	var formErrors []string
	if contactName == "" {
		formErrors = append(formErrors, "Your name cannot be empty.")
	}
	if contactEmail == "" {
		formErrors = append(formErrors, "You must provide a valid email address.")
	}
	if contactMessage == "" {
		formErrors = append(formErrors, "You must provide message.")
	}

	inst.contactFormSessionsLock.Lock()
	defer inst.contactFormSessionsLock.Unlock()
	contactSession, exists := inst.contactFormSessions[contactSessionToken]
	if !exists {
		formErrors = append(formErrors, "Invalid session token.")
		return false, formErrors, formData
	}

	// Delete session if it exists so it can't be re-used.
	delete(inst.contactFormSessions, contactSessionToken)

	contactSumAnswerInt, err := strconv.Atoi(contactSumAnswer)
	if err != nil {
		formErrors = append(formErrors, "Invalid answer for the two number sum.")
		return false, formErrors, formData
	}
	if contactSumAnswerInt != contactSession.sumAnswer {
		formErrors = append(formErrors, "Invalid answer for the two number sum.")
		return false, formErrors, formData
	}

	if time.Since(contactSession.createdAt).Seconds() < 4 {
		formErrors = append(formErrors, "Form failed to validate. Please wait and try again.")
	}

	if len(formErrors) > 0 {
		return false, formErrors, formData
	}

	return true, nil, formData
}
