package data

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidPriceFormat = errors.New("invalid price format")

type Price float64

func (p Price) MarshalJSON() ([]byte, error) {
	// jsonValue := fmt.Sprintf("$%d", p) from int32 to float64
	jsonValue := "$ " + strconv.FormatFloat(float64(p), 'f', -1, 64)
	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil
}

func (p *Price) UnmarshalJSON(jsonValue []byte) error {
	// We expect that the incoming JSON value will be a string in the format
	// "<runtime> mins", and the first thing we need to do is remove the surrounding
	// double-quotes from this string. If we can't unquote it, then we return the
	// ErrInvalidRuntimeFormat error.
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidPriceFormat
	}
	// Split the string to isolate the part containing the number.
	parts := strings.Split(unquotedJSONValue, " ")
	// Sanity check the parts of the string to make sure it was in the expected format.
	// If it isn't, we return the ErrInvalidRuntimeFormat error again.
	if len(parts) != 2 || parts[0] != "$" {
		return ErrInvalidPriceFormat
	}
	// Otherwise, parse the string containing the number into an int32. Again, if this
	// fails return the ErrInvalidRuntimeFormat error.
	i, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return ErrInvalidPriceFormat
	}
	// Convert the int32 to a Runtime type and assign this to the receiver. Note that we
	// use the * operator to deference the receiver (which is a pointer to a Runtime
	// type) in order to set the underlying value of the pointer.
	*p = Price(i)
	return nil
}
