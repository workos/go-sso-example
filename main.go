package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/workos-inc/workos-go/pkg/sso"
)

func main() {
	var conf struct {
		Addr        string
		APIKey      string
		ProjectID   string
		RedirectURI string
		Domain      string
	}
	flag.StringVar(&conf.Addr, "addr", ":3042", "The server addr.")
	flag.StringVar(&conf.APIKey, "api-key", "", "The WorkOS API key.")
	flag.StringVar(&conf.ProjectID, "project-id", "", "The WorkOS project id.")
	flag.StringVar(&conf.RedirectURI, "redirect-uri", "", "The redirect uri.")
	flag.StringVar(&conf.Domain, "domain", "", "The domain used to register a WorkOS SSO connection.")
	flag.Parse()

	log.Printf("launching sso demo with configuration: %+v", conf)

	// Configure the WorkOS SSO SDK:
	sso.Configure(conf.APIKey, conf.ProjectID, conf.RedirectURI)

	// Handle login:
	http.Handle("/login", sso.Login(sso.GetAuthorizationURLOptions{
		Domain: conf.Domain,
	}))

	// Handle login redirect:
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("callback is called with %s", r.URL)

		// Retrieving user profile:
		profile, err := sso.GetProfile(context.Background(), sso.GetProfileOptions{
			Code: r.URL.Query().Get("code"),
		})
		if err != nil {
			log.Printf("get profile failed: %s", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// Display user profile:
		b, err := json.MarshalIndent(profile, "", "    ")
		if err != nil {
			log.Printf("encoding profile failed: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		w.Write(b)

		log.Printf("user is logged with profile: %s", b)
	})

	if err := http.ListenAndServe(conf.Addr, nil); err != nil {
		log.Panic(err)
	}
}
