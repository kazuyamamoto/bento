package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Google App Engine がポートを指定する
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)

	s := &slackHandler{
		shouldNotify: newWeekday().isToday,
		// TODO Incoming Webhook の URL をセット
		url: "https://hooks.slack.com/services/XXXXXXXX/YYYYYYYYY/ZZZZZZZZZZZZZZZZZZZZZZZZ",
		vendors: []*vendor{
			newTamagoya(),
			newAzuma2020(),
		},
	}

	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
