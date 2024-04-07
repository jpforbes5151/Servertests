package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jpforbes5151/Servertests/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		err := db.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	//creates a new server
	mux := http.NewServeMux()

	// points middleware towards the appropriate filepath
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	// handles if the Server is ready to accept a request
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	// resets the metric for how many times a resource has been hit
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)

	//handles how to create a user
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	//handles how to login a user
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	// provides metrics on how many times an endpoint has been hit
	mux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)

	// handles the creation of a Chirp
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)

	// handles the retreival of a Chirp
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)

	// retrieves a specific Chirp
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)

	//provides admin metrics on how many times a resource has been hit
	mux.HandleFunc("GET /admin/metrics", apiCfg.adminHandlerMetrics)

	// initiates security protocols for headers
	corsMux := middlewareCors(mux)

	//server configuration
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
