package producer

import (
	"workshop.4.2/internal/model"
)

func Orders() <-chan model.ClientID {
	result := make(chan model.ClientID, 5)
	go func() {
		defer close(result)
		for i := 0; i < 5; i++ {
			result <- model.ClientUserID
		}
	}()
	return result
}
