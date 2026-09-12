package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func BoolPtr(v bool) *bool {
	return &v
}

func StringPtr(s string) *string {
	return &s
}

type Envelope map[string]any

func WriteJson(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for k, v := range headers {
		for _, header := range v {
			w.Header().Add(k, header)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func ReadJson(w http.ResponseWriter, r *http.Request, dst any) error {

	//Set max req body to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		//Start triage
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON at character %d", syntaxError.Offset)

		//Decode() may return io.ErrUnexpectedEOF error for syntax errs in JSON
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		//Occurs when JSON value is wrong type for destination. Include specific field if available for easier debugging
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("Body contains incorrect JSON for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("Body contains badly-formed JSON at character %d", unmarshalTypeError.Offset)

		//Occurs when body is empty
		case errors.Is(err, io.EOF):
			return errors.New("Body must not be empty")

		//Occurs when JSON contains field cannot be mapped to target destination
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains uknown key %s", fieldName)

		//Occurs when message body has exceeded max size
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		//Occurs when we pass a non-nil pointer as the decode destination
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}
	return nil
}
