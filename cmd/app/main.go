package main

import (
	"encoding/json"
	"fmt"
	"git.perkbox.io/backend-services/intercom-search/client/helpdocs"
	"git.perkbox.io/backend-services/intercom-search/client/intercom"
	"github.com/bykovme/gotrans"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	_ "net/http/pprof"
	"os"
)

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

	// Intercom will call this first to initialise some key-value pair. We don't need to set any yet
	// Intercom API doesn't accept an empty results array, so we're setting a `key` and `value` anyway.
	r.Post("/configure", func(w http.ResponseWriter, r *http.Request){
		results := intercom.ConfigResults{
			map[string]string{"key": "value"},
		}

		o, _ := json.Marshal(results)
		log.Println(string(o))
		_, err = w.Write(o)

	})

	// Second call by the Intercom will made to this endpoint and expect a canvas component to draw
	// on the chat display.
	r.Post("/initialize", func(w http.ResponseWriter, r *http.Request){
		b, _ := httputil.DumpRequest(r, true)
		log.Println(string(b))

		canvas := intercom.NewCanvas()
		input := intercom.NewInput(gotrans.Tr("en", "search"), gotrans.Tr("en", "search_for_articles"))

		canvas.AddComponent(input)

		o, _ := json.Marshal(canvas)
		_, err = w.Write(o)
		return

	})

	// After initialisation, Intercom will post the request to this endpoint and this will become the primary
	// endpoint for the rest of the flow.
	r.Post("/search", func(w http.ResponseWriter, r *http.Request){
		b, _ := httputil.DumpRequest(r, true)
		log.Println(string(b))

		request := &intercom.Request{}
		err := json.NewDecoder(r.Body).Decode(request)

		if err != nil {
			log.Println("Can't decode intercom request into intercom.Request{}", err)
		}

		fmt.Printf("%+v", request)

		// get the region to determine the language
		region := request.User.CustomAttributes.Region
		log.Println("Region is : ", region)

		// set the locale. it's either English or French and hardcoded because we can't rely on the language param
		// when testing from UK. can be updated later.
		locale := "en"

		if region == "fr-FR" {
			locale = "fr"
		}

		// prepare the canvas and add the search box to it
		canvas := intercom.NewCanvas()
		input := intercom.NewInput(gotrans.Tr(locale, "search"), gotrans.Tr(locale, "search_for_articles"))


		canvas.AddComponent(input)

		// check if a search request has been sent and do the search if it has
		if query, ok := request.InputValues["query"]; ok {

			log.Println("QUERY IS : ", query)
			log.Println(request.User.CustomAttributes)

			// we have 3 different accounts, thus 3 API keys
			locations := map[string]string{
				"en-GB": os.Getenv("HELPDOCS_API_KEY_GB"),
				"fr-FR": os.Getenv("HELPDOCS_API_KEY_FR"),
				"en-AU": os.Getenv("HELPDOCS_API_KEY_AU"),
			}

			log.Println(locations)

			// and update it if another location is provided
			if val, ok := locations[region]; ok {
				log.Println("Selected locale is ", region)
				hd.SetAuthToken(val)
			} else {
				log.Println("Selected locale has not found. Switching to default locale, which is en-GB")
				hd.SetAuthToken(locations["en-GB"])
			}

			response, err := hd.Search(query)

			if err != nil {
				log.Println("An error occurred while calling Search :", err)
			}

			log.Println(response)

			// prepare the list response if the number of articles is greater than zero
			if len(response.Articles) > 0 {

				component := intercom.ListResponse{}
				component.Type = "list"
				component.Disabled = false

				for _, item := range response.Articles {
					component.Items = append(component.Items, intercom.ListResponseItem{
						"item",
						item.ArticleId,
						item.Title,
						item.Description,
						"",
						0,
						0,
						false,
						false,
						intercom.SheetAction{"sheet", "https://61958bdd.ngrok.io/show/" + item.ArticleId},
					})

					log.Println("Item URL is : ", item.Url)

					canvas.AddComponent(component)
				}

			}

		}

		o, _ := json.Marshal(canvas)
		_, err = w.Write(o)
		return

	})

	r.Post("/show/{articleId}", func(w http.ResponseWriter, r *http.Request){
		b, _ := httputil.DumpRequest(r, true)
		log.Println(string(b))

		//hd.SetAuthToken(os.Getenv("HELPDOCS_API_KEY_FR"))

		articleId := chi.URLParam(r, "articleId")
		article, err := hd.GetArticle(articleId)

		if err != nil {
			log.Println("Error while fetching the article : ", err)
		}

		log.Println(article)

		temp, err := template.New("article").Parse(`
			<html>
			<head>
			<title>{{.Title}}</title>
			</head>
			<body>
			<h1>{{.Title}}</h1>
			<div>
			{{.Body}}
			</div>
			</body>
			</html>
		`)

		if err != nil {
			log.Println("Failed to create the template : ", err)
		}

		err = temp.Execute(w, article.Article)

		if err != nil {
			log.Println("Failed to print the template : ", err)
		}

	})

	err = http.ListenAndServe("127.0.0.1:3000", r)

	if err != nil {
		panic(err)
	}

}
