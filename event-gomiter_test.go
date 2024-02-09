package eventgomiter

import (
	"fmt"
	"testing"
)

var emitter EventEmitter = NewChannelEventEmitter()

func TestRegisterHandler(t *testing.T) {
	var name EventIdentifier = "test"
	handler := func(e Event) {
		fmt.Println(e)
	}

	err := emitter.RegisterHandler(name, handler)
	if err != nil {
		t.Fatalf("Emitter does not register handlers properly")
	}
}
