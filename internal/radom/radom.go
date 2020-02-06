package radom

import (
	"time"
	"math/rand"
)


func Random(arr []string) {
	if len(arr) <= 0 {
		return
	}
	rand.Seed(time.Now().Unix())
	for i := len(arr) - 1; i >= 0; i-- {
		num := rand.Intn(len(arr))
		arr[i], arr[num] = arr[num], arr[i]
	}

	return

	
}
