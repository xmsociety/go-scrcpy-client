package test

import (
	"fmt"
	"testing"
)

func TestEquaBuf(t *testing.T) {
	fmt.Println([]byte("\x00")[0] == []byte("\x00")[0])
}
