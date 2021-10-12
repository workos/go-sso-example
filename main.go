package main

import (
	"context"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/workos-inc/workos-go/pkg/sso"
)

func main() {
	var conf struct {
		Addr        string
		APIKey      string
		ClientID    string
		RedirectURI string
		Domain      string
		Provider    string
	}

	type Profile struct {
		ProfileData string
	}

	flag.StringVar(&conf.Addr, "addr", ":3042", "The server addr.")
	flag.StringVar(&conf.APIKey, "api-key", "sk_test_a2V5XzAxRkExMkM3TTNSTldFNUNKSEFNUUVZQ1pTLDJtb3drUExOTk9vT3dDc1NDRTZnRUVVQ28", "The WorkOS API key.")
	flag.StringVar(&conf.ClientID, "client-id", "client_01FA12C7QV793K318T2G1V3E7X", "The WorkOS project id.")
	flag.StringVar(&conf.RedirectURI, "redirect-uri", "http://localhost:3042/callback", "The redirect uri.")
	flag.StringVar(&conf.Domain, "domain", "gmail.com", "The domain used to register a WorkOS SSO connection.")
	flag.StringVar(&conf.Provider, "provider", "MicrosoftOAuth", "The OAuth provider used for the SSO connection.")
	flag.Parse()

	log.Printf("launching sso demo with configuration: %+v", conf)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Configure the WorkOS SSO SDK:
	sso.Configure(conf.APIKey, conf.ClientID)

	// Handle login
	http.Handle("/login", sso.Login(sso.GetAuthorizationURLOptions{
		//Instead of domain, you can now use connection ID to associate a user to the appropriate connection.
		Domain:      conf.Domain,
		RedirectURI: conf.RedirectURI,
	}))

	// Handle login redirect:
	tmpl := template.Must(template.ParseFiles("./static/logged_in.html"))
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("callback is called with %s", r.URL)

		// Retrieving user profile:
		profile, err := sso.GetProfileAndToken(context.Background(), sso.GetProfileAndTokenOptions{
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
		print(b)
		if err != nil {
			log.Printf("encoding profile failed: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		// Convert the profile to a string
		stringB := string(b)
		data := Profile{
			ProfileData: stringB,
		}

		// Render the template
		tmpl.Execute(w, data)
	})

	if err := http.ListenAndServe(conf.Addr, nil); err != nil {
		log.Panic(err)
	}
}
