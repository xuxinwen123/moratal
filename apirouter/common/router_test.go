package common

import (
	"fmt"
	"testing"
)

func TestParsePattern(t *testing.T) {
	pattern := ParsePattern("/api/user/{id:int}")
	fmt.Println(pattern)
}
