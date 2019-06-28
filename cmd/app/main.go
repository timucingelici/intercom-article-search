package main

import (
	"github.com/bykovme/gotrans"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"github.com/timucingelici/intercom-article-search/client/helpdocs"
	"github.com/timucingelici/intercom-article-search/internal"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

// TODO: add tests
func main() {

	// memory profiling
	go func() {
		log.Println(http.ListenAndServe("127.0.0.1:9999", nil))
	}()

	// load the env vars
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// initialise the locales
	if err := gotrans.InitLocales("internal/localizations"); err != nil {
		log.Fatal("Error while loading the locales : ", err)
	}

	// routing
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// init helpdocs client
	hd, err := helpdocs.NewClient(http.DefaultClient, "https://api.helpdocs.io/v1", "")

	internal.SetUpRoutes(r, hd)

	err = http.ListenAndServe("127.0.0.1:"+os.Getenv("SERVICE_PORT"), r)

	if err != nil {
		panic(err)
	}

}
