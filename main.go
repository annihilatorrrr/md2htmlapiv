package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Frameit struct {
	Type   string `json:"type"`
	Button string `json:"button"`
	Text   string `json:"text"`
}

func mdfunc(w http.ResponseWriter, r *http.Request) {
	log.Print("Endpoint main Triggered!")
	htmlText, btns := MD2HTMLButtonsV2(r.FormValue("rtext"))
	formatedtrext := Frameit{
		Type:   "success!",
		Button: fmt.Sprintf("%v", btns),
		Text:   htmlText,
	}
	err := json.NewEncoder(w).Encode(formatedtrext)
	if err != nil {
		log.Print(err.Error())
	}
}

func homePage(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprint(w, "Welcome! API_V: v1; Thanks to @PaulSonOfLars, @AmarnathCJD @Divkix @anonyindian !")
	if err != nil {
		log.Println(err.Error())
	}
	log.Print("Endpoint Hit: HomePage!")
}

// Main function
func main() {
	// Init the mux router
	router := mux.NewRouter().StrictSlash(true)
	port := os.Getenv("PORT")
	// serve it
	if port == "" {
		port = "80"
	}
	server := &http.Server{
		Addr:         "0.0.0.0:" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 3 * time.Second,
		IdleTimeout:  6 * time.Second,
	}
	// Route handles & endpoints
	router.HandleFunc("/", homePage)
	router.HandleFunc("/md2htmlbv2/", mdfunc).Methods("POST")
	log.Printf("Started Api at %s !", os.Getenv("LINK"))
	log.Fatal(server.ListenAndServe())
}
