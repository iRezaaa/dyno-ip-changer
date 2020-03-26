package main

import "fmt"

type CustomError struct {
	Type    ErrorType
	Details string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("code : %d : message %s", e.Type, e.Details)
}
