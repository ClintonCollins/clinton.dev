package web

import (
	"html/template"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"clinton.dev/internal/email"
	"clinton.dev/internal/utils"
	"github.com/go-chi/chi/v5"
)

type instance struct {
	frontendPath            string
	indexTemplate           *template.Template
	contactTemplate         *template.Template
	httpServer              *http.Server
	email                   *email.Instance
	enableStaticFiles       bool
	contactEmailFrom        string
	contactEmailTo          string
	contactFormSessions     map[string]contactFormSession
	contactFormSessionsLock *sync.RWMutex
}

// contactFormSession will be created every time the contact form is loaded. It's basically a way of storing sessions.
// It's how we do spam prevention on the contact form. This could be replaced in the future if it becomes a memory hog.
// Could be replaced with signed input fields via HMAC.
type contactFormSession struct {
	createdAt time.Time
	sumAnswer int // Answer for provided simple math problem.
}

// Generate a random seed on load.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// New creates a new struct with methods to start the public web server. It takes http.Server so you can configure settings the way you want.
// It also parses all the necessary template files and gets the frontend path for you.
func New(httpServer *http.Server, emailInstance *email.Instance, enableStaticFileServer bool,
	contactEmailFrom, contactEmailTo string,
) (*instance, error) {
	p, err := utils.GetRelativeDirectory("/frontend")
	if err != nil {
		return nil, err
	}
	inst := &instance{
		frontendPath:            p,
		indexTemplate:           nil,
		httpServer:              httpServer,
		email:                   emailInstance,
		enableStaticFiles:       enableStaticFileServer,
		contactFormSessions:     make(map[string]contactFormSession),
		contactFormSessionsLock: &sync.RWMutex{},
		contactEmailTo:          contactEmailTo,
		contactEmailFrom:        contactEmailFrom,
	}
	err = inst.parseTemplates()
	if err != nil {
		return nil, err
	}
	go inst.cleanupOldContactSessionsRunner()
	return inst, nil
}

// Listen setups a router and then blocks listening for http connections.
func (inst *instance) Listen() error {

	r := chi.NewRouter()
	// Serve static files directly instead of letting a webserver like Nginx handle them.
	if inst.enableStaticFiles {
		staticFiles := http.FileServer(http.Dir(inst.frontendPath + "/static"))
		mediaFiles := http.FileServer(http.Dir(inst.frontendPath + "/media"))

		r.Handle("/static/*", http.StripPrefix("/static", staticFiles))
		r.Handle("/media/*", http.StripPrefix("/media", mediaFiles))
	}

	// Public routes
	r.HandleFunc("/", inst.getIndex)
	r.HandleFunc("/contact", inst.contact)

	// Public API
	r.HandleFunc("/api/contact", inst.apiContact)

	inst.httpServer.Handler = r

	return inst.httpServer.ListenAndServe()
}
