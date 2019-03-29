# Problem

Perkbox now has 3 different regions, UK, FR and AU yet Intercom's knowledge base doesn't support multiple categories or multi-region knowledge base.

# Solution 

This app acts a bridge between Intercom's messenger and a 3rd party knowledge base and fetches the right content by the user's region.

# Setup

* Clone the repo
* Rename/copy `.env-example` as `.env`
* Install the dependencies with `dep ensure`
* Build & run `cmd/app/main.go`

# Notes

* Entry point for the app is `cmd/app/main.go`
* Tests to be added