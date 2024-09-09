package handler

import "os"

func getAllowBrowsing() bool {
	return os.Getenv("ALLOW_BROWSING") == "true"
}
