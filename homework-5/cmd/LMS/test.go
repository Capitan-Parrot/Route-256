package main

import (
	"fmt"
	"gitlab.ozon.dev/homework5/internal/pkg/server"
	"reflect"
)

func main() {
	t := reflect.TypeOf(&server.Transport{})

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		fmt.Println(m.Name)
	}
}
