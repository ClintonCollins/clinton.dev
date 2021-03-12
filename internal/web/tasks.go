package web

import (
	"time"
)

// cleanupOldContactSessions loops through all the contact sessions and deleted anything over an hour old.
func (inst *instance) cleanupOldContactSessions() {
	inst.contactFormSessionsLock.Lock()
	for token, s := range inst.contactFormSessions {
		if time.Since(s.createdAt).Hours() > 1 {
			delete(inst.contactFormSessions, token)
		}
	}
	inst.contactFormSessionsLock.Unlock()
}

// cleanupOldContactSessionsRunner is meant to be run in a go routine to periodically trigger the session cleanup
// function every 5 minutes.
func (inst *instance) cleanupOldContactSessionsRunner() {
	throttle := time.NewTicker(time.Minute * 5)
	defer throttle.Stop()
	for {
		<- throttle.C
		inst.cleanupOldContactSessions()
	}
}
