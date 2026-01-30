package store

import "fmt"

func generateStoreUploadKey(storeID int64) string {
	return fmt.Sprintf("stores/%d/site.json", storeID)
}
