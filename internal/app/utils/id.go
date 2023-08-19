package utils

import (
	"fmt"
	"strconv"
)

func ParseId(id string) int {
	intId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Print("Unable to parse id")
	}

	return intId
}
