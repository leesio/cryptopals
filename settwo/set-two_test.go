package settwo

import (
	"fmt"
	"testing"
)

func TestPadding(t *testing.T) {
	r := pad(20, []byte("YELLOW SUBMARINE"))
	fmt.Println(r)
}
