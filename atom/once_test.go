package atom

import (
	"fmt"
	"testing"
)

func TestOnce(t *testing.T) {
	Once(func() {
		fmt.Println("v")
	})
}
