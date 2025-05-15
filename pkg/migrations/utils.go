package migrations

import (
	"fmt"
	"regexp"
)

func makeMigrateUrl(dbUrl string) string {
	urlRe := regexp.MustCompile(`^[^\\?]+`)
	url := urlRe.FindString(dbUrl)

	sslModeRe := regexp.MustCompile("(sslmode=)[a-zA-Z0-9]+")
	sslMode := sslModeRe.FindString(dbUrl)

	return fmt.Sprintf("%s?%s", url, sslMode)
}
