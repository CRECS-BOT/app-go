package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

// handler per /api (per ora dummy)
func apiHandler(w http.ResponseWriter, r *http.Request) {
	// esempio: semplice JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok","message":"API alive OK ;)"}`))
}

// handler per /telegram (dummy, solo 200 OK)
func telegramHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func main() {
	// flag --port, default 8080
	port := flag.Int("port", 8080, "Porta HTTP su cui avviare il server")
	flag.Parse()

	mux := http.NewServeMux()

	// --- ROUTE API ---
	// /api oppure /api/... (puoi anche usare /api/ se vuoi un sotto-router)
	mux.HandleFunc("/api", apiHandler)
	// se vuoi che /api/... vada uguale:
	// mux.HandleFunc("/api/", apiHandler)

	// --- ROUTA TELEGRAM WEBHOOK ---
	mux.HandleFunc("/telegram", telegramHandler)
	// o /telegram/webhook se vuoi qualcosa di più specifico
	// mux.HandleFunc("/telegram/webhook", telegramHandler)

	// --- STATICI WEBSITE su / ---
	// Tutto ciò che non matcha sopra e parte da / va su website
	websiteDir := http.Dir("website")
	websiteFS := http.FileServer(websiteDir)
	mux.Handle("/", websiteFS)

	// --- STATICI WEBAPP su /webapp ---
	// Serve i file dentro ./webapp quando la URL inizia con /webapp
	webappDir := http.Dir("webapp")
	webappFS := http.StripPrefix("/webapp", http.FileServer(webappDir))
	mux.Handle("/webapp/", webappFS)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Server starting on %s\n", addr)

	// controllo banale che le cartelle esistano (solo log, non blocco)
	if _, err := os.Stat("website"); os.IsNotExist(err) {
		log.Println("WARNING: cartella 'website' non trovata")
	}
	if _, err := os.Stat("webapp"); os.IsNotExist(err) {
		log.Println("WARNING: cartella 'webapp' non trovata")
	}

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
