package helpdocs

import (
	"os"
	"strings"
)


func GetApiKeys() map[string]string {
	return map[string]string{
		"en-GB": os.Getenv("HELPDOCS_API_KEY_GB"),
		"fr-FR": os.Getenv("HELPDOCS_API_KEY_FR"),
		"en-AU": os.Getenv("HELPDOCS_API_KEY_AU"),
	}
}

func GetApiKeyByRegion(region string) string {
	keys := GetApiKeys()
	if val, ok := keys[region]; ok {
		return val
	}

	return keys["en-GB"]
}

func GetLocaleByRegion(region string) string {
	s := strings.Split(region, "-")
	return s[0]
}

