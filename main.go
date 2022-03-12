package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/gorilla/mux"
)

type Frameit struct {
	Type   string `json:"type"`
	Button string `json:"button"`
	Text   string `json:"text"`
}

func mdfunc(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint main Triggered! Received:", r.FormValue("rtext"))
	htmlText, btns := MD2HTMLButtonsV2(r.FormValue("rtext"))
	formatedtrext := Frameit{
		Type:   "success!",
		Button: fmt.Sprintf("%v", btns),
		Text:   htmlText,
	}
	err := json.NewEncoder(w).Encode(formatedtrext)
	if err != nil {
		log.Println(err.Error())
	}
}

func homePage(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "Welcome! API_V: v1; Thanks to @PaulSonOfLars, @AmarnathCJD @Divkix @anonyindian !")
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Endpoint Hit: HomePage!")
}

// Main function
func main() {
	debug.SetGCPercent(3)
	// Init the mux router
	router := mux.NewRouter().StrictSlash(true)
	port := os.Getenv("PORT")
	// serve it
	log.Println("API Link:", os.Getenv("LINK"))
	log.Println("Server at", port)
	server := &http.Server{
		Addr:         "0.0.0.0:" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 3 * time.Second,
		IdleTimeout:  7 * time.Second,
	}
	// Route handles & endpoints
	router.HandleFunc("/", homePage)
	router.HandleFunc("/md2htmlbv2/", mdfunc).Methods("POST")
	runtime.GC()
	debug.FreeOSMemory()
	log.Fatal(server.ListenAndServe())
}
