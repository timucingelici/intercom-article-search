package internal

import (
	"encoding/json"
	"git.perkbox.io/poc/intercom-article-search/client/helpdocs"
	"git.perkbox.io/poc/intercom-article-search/client/intercom"
	"github.com/bykovme/gotrans"
	"github.com/go-chi/chi"
	"html/template"
	"log"
	"net/http"
	"os"
)

func configureHandler(w http.ResponseWriter, r *http.Request) {
	results := intercom.ConfigResults{
		map[string]string{"key": "value"},
	}

	sendResponse(results, w)
	return
}

func initializeHandler(w http.ResponseWriter, r *http.Request) {
	canvas := intercom.NewCanvas()
	input := intercom.NewInput(gotrans.Tr("en", "search"), gotrans.Tr("en", "search_for_articles"))

	canvas.AddComponent(input)

	sendResponse(canvas, w)
	return

}

func searchHandler(hd *helpdocs.Helpdocs) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {

		request := &intercom.Request{}
		err := json.NewDecoder(r.Body).Decode(request)

		if err != nil {
			log.Println("Failed to decode intercom request into the struct : ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// get the region to determine the language
		region := request.User.CustomAttributes.Region

		// set the locale. it's either English or French and hardcoded because we can't rely on the language param
		// when testing from UK. can be updated later.
		locale := helpdocs.GetLocaleByRegion(region)

		// prepare the canvas and add the search box to it
		canvas := intercom.NewCanvas()
		input := intercom.NewInput(gotrans.Tr(locale, "search"), gotrans.Tr(locale, "search_for_articles"))

		canvas.AddComponent(input)

		// check if a search param has been sent
		query, ok := request.InputValues["query"]

		if !ok {
			sendResponse(canvas, w)
			return
		}

		hd.SetAuthToken(helpdocs.GetApiKeyByRegion(region))

		results, err := hd.Search(query)

		if err != nil {
			log.Println("An error occurred while calling Search :", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(results.Articles) > 0 {

			component := intercom.ListResponse{}
			component.Type = "list"
			component.Disabled = false

			for _, item := range results.Articles {

				component.Items = append(component.Items, intercom.ListResponseItem{
					Type:         "item",
					Id:           item.ArticleId,
					Title:        item.Title,
					Subtitle:     item.Description,
					Image:        "",
					ImageWidth:   0,
					ImageHeight:  0,
					RoundedImage: false,
					Disabled:     false,
					Action:       intercom.SheetAction{"sheet", os.Getenv("SERVICE_URL") + "/show/" + item.ArticleId + "?region=" + region},
				})

				canvas.AddComponent(component)
			}

		}

		sendResponse(canvas, w)
	}

	return http.HandlerFunc(fn)
}

func showArticleHandler(hd *helpdocs.Helpdocs) http.HandlerFunc {

	fn := func(w http.ResponseWriter, r *http.Request) {

		//intercom doesn't send the region data back, so we'll extract it from the URL
		err := r.ParseForm()

		if err != nil {
			log.Println("Error while parsing from data : ", err)
		}

		region := r.Form.Get("region")

		hd.SetAuthToken(helpdocs.GetApiKeyByRegion(region))

		articleId := chi.URLParam(r, "articleId")
		article, err := hd.GetArticle(articleId)

		if err != nil {
			log.Println("Error while fetching the article : ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Println(article)

		tempVars := map[string]interface{}{
			"Title": article.Article.Title,
			"Body":  template.HTML(article.Article.Body),
		}

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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = temp.Execute(w, tempVars)

		if err != nil {
			log.Println("Failed to print the template : ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	return http.HandlerFunc(fn)
}
