package stock

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	res, err := Get([]string{"sh601360", "sz000555"})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", res)
}
