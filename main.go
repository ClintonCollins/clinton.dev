package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"clinton.dev/internal/email"
	"clinton.dev/internal/utils"
	"clinton.dev/internal/web"
	"github.com/BurntSushi/toml"
)

type configuration struct {
	PostmarkServerToken  string `toml:"PostmarkServerToken"`
	PostmarkAccountToken string `toml:"PostmarkAccountToken"`
	WebServerIP          string `toml:"WebServerIP"`
	WebServerPort        int    `toml:"WebServerPort"`
	ServeStaticFiles     bool   `toml:"ServeStaticFiles"`
	ContactFormFromEmail string `toml:"ContactFormFromEmail"`
	ContactFormToEmail   string `toml:"ContactFormToEmail"`
}

func getConfiguration() (*configuration, error) {
	configPath, err := utils.GetRelativeDirectory("config.toml")
	if err != nil {
		return nil, err
	}
	config := &configuration{}
	_, err = toml.DecodeFile(configPath, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func main() {
	config, err := getConfiguration()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error reading configuration file. ERROR: %v", err))
	}

	httpServer := &http.Server{
		Addr:              net.JoinHostPort(config.WebServerIP, strconv.Itoa(config.WebServerPort)),
		ReadTimeout:       time.Second * 15,
		ReadHeaderTimeout: time.Second * 15,
		WriteTimeout:      time.Second * 15,
		IdleTimeout:       time.Minute * 60,
	}

	// Configure and get postmark email methods.
	postmarkEmailInstance := email.New(config.PostmarkServerToken, config.PostmarkAccountToken)

	// Configure and get access to webserver methods.
	webServer, err := web.New(httpServer, postmarkEmailInstance, config.ServeStaticFiles,
		config.ContactFormFromEmail, config.ContactFormToEmail)
	if err != nil {
		log.Fatal(err)
	}

	// Self explanatory
	err = webServer.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
