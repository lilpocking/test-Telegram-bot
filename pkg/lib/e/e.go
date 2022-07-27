package e

import "fmt"

// Wrap func that just wrapped message and err without checking, and return it's in error type.
func Wrap(message string, err error) error {
	return fmt.Errorf("%s: %w", message, err)
}

/*

	WrapIfErr func that check err parametr is nil
If it's nil than return nil. Else return wrapped message in error type
*/
func WrapIfErr(message string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %w", message, err)
	}
	return nil
}
