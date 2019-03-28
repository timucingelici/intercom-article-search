package helpdocs

import (
	"os"
	"strings"
)

type apiKeys struct {}

func (k apiKeys) All() map[string]string {
	return map[string]string {
		"en-GB": os.Getenv("HELPDOCS_API_KEY_GB"),
		"fr-FR": os.Getenv("HELPDOCS_API_KEY_FR"),
		"en-AU": os.Getenv("HELPDOCS_API_KEY_AU"),
	}
}

func (k apiKeys) Get(region string) string {
	keys := k.All()
	if val, ok := keys[region]; ok {
		return val
	}

	return keys["en-GB"]
}

func (k apiKeys) GetLocale(region string) string {
	r := k.Get(region)
	s := strings.Split(r, "-")
	return s[0]
}

func ApiKeys() apiKeys {
	return apiKeys{}
}