package helpdocs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Helpdocs struct {
	baseUrl   *url.URL
	userAgent string
	client    *http.Client
	authToken string
}

func NewClient(client *http.Client, baseUrl string, userAgent string) (*Helpdocs, error) {

	if client == nil {
		client = http.DefaultClient
	}

	u, err := url.Parse(baseUrl)

	if err != nil {
		return nil, err
	}

	if userAgent == "" {
		userAgent = "Go Client for HelpDocs API"
	}

	return &Helpdocs{
		u,
		userAgent,
		client,
		"",
	}, nil
}

func (hd *Helpdocs) SetAuthToken(token string) {
	hd.authToken = token
}

func (hd *Helpdocs) newRequest(method, path string, body interface{}) (*http.Request, error) {

	rel := &url.URL{Path: path}
	u := hd.baseUrl.ResolveReference(rel)

	var buf io.ReadWriter

	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)

		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)

	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", hd.userAgent)

	if hd.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+hd.authToken)
	}

	return req, nil
}

func (hd *Helpdocs) call(method string, endpoint string, queryParams map[string]string, response interface{}, body interface{}) error {

	// create the url
	u := hd.baseUrl.Path + "/" + endpoint

	// create the request
	req, err := hd.newRequest(method, u, nil)

	if err != nil {
		return err
	}

	// add query params to the request
	q := req.URL.Query()

	for k, v := range queryParams {
		q.Add(k, v)
	}

	q.Add("query", "cinema")

	req.URL.RawQuery = q.Encode()

	// send the request
	res, err := hd.client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	// parse response into the struct
	err = json.NewDecoder(res.Body).Decode(response)

	if err != nil {
		log.Println(err)
	}

	return err
}

func (hd *Helpdocs) Search(q string) (*Search, error) {

	var data Search

	queryStrings := map[string]string{"query": q}

	if err := hd.call("GET", "search", queryStrings, &data, nil); err != nil {
		return nil, err
	}

	fmt.Println(&data)

	return &data, nil
}

func (hd *Helpdocs) GetArticle(articleId string) (*ArticleResponse, error) {

	var data ArticleResponse

	if err := hd.call("GET", "article/"+articleId, nil, &data, nil); err != nil {
		return nil, err
	}

	return &data, nil
}