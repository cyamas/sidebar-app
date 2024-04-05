package app

import (
	"math/rand"
	"time"
)

func generateUniqueID(idList *[]int) int {
	rand.New(rand.NewSource(time.Now().Unix()))
	candidate := rand.Intn(900000) + 100000
	if IsUniqueID(candidate, idList) {
		return candidate
	} else {
		return generateUniqueID(idList)
	}
}

func IsUniqueID(candidate int, idList *[]int) bool {
	if len(*idList) > 0 {
		for _, id := range *idList {
			if candidate == id {
				return false
			}
		}
	}
	return true
}
