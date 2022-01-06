package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type Frameit struct {
	Type   string `json:"type"`
	Button string `json:"button"`
	Text   string `json:"text"`
}

func mdfunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint main Triggred! Recived:", r.FormValue("rtext"))
	htmlText, btns := MD2HTMLButtonsV2(r.FormValue("rtext"))
	var formatedtrext = Frameit{
		Type:   "success!",
		Button: fmt.Sprintf("%v", btns),
		Text:   htmlText,
	}
	json.NewEncoder(w).Encode(formatedtrext)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome! API_V: v1; Thanks to @PaulSonOfLars, @AmarnathCJD @Divkix @anonyindian !")
	fmt.Println("Endpoint Hit: HomePage!")
}

// Main function
func main() {

	// Init the mux router
	router := mux.NewRouter().StrictSlash(true)

	// Route handles & endpoints
	router.HandleFunc("/", homePage)

	// call it
	router.HandleFunc("/md2htmlbv2/", mdfunc).Methods("POST")
	port := os.Getenv("PORT")
	// serve it
	fmt.Println("API Link:", os.Getenv("LINK"))
	fmt.Println("Server at", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
