package internal

import (
	"github.com/go-chi/chi"
	"github.com/timucingelici/intercom-article-search/client/helpdocs"
)

func SetUpRoutes(r *chi.Mux, hd *helpdocs.Helpdocs) {

	// Intercom will call this first to initialise some key-value pair. We don't need to set any yet
	// Intercom API doesn't accept an empty results array, so we're setting a `key` and `value` anyway.
	r.Post("/configure", configureHandler)

	// Second call by the Intercom will made to this endpoint and expect a canvas component to draw
	// on the chat display.
	r.Post("/initialize", initializeHandler)

	// After initialisation, Intercom will post the request to this endpoint and this will become the primary
	// endpoint for the rest of the flow.
	r.Post("/search", searchHandler(hd))

	// Finally, when a user click on a result, this route will be triggered and fetch the article and print it
	// in a html template.
	r.Post("/show/{articleId}", showArticleHandler(hd))
}
