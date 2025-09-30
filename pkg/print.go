package pkg

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	Name string `json:"event,omitempty"`
	Data any    `json:"data,omitempty"`
}

func MustPrint(e Event) {
	data, err := json.MarshalIndent(e, " ", "   ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n%s\n", string(data))
}
