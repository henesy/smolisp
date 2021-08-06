package main

import (
	"fmt"
	"errors"
)


// Generic error message builder
func e(v ...interface{}) error {
	if len(v) < 2 {
		return errors.New(fmt.Sprint(v[0]))
	}
	return errors.New(fmt.Sprint(v...))
}

