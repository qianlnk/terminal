package terminal

import (
	"fmt"
	"testing"
)

func TestGetch(t *testing.T) {
	for {
		fmt.Printf("> ")
		getch()
	}
}
