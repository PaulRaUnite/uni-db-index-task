package urlval

import (
	"net/url"
	"reflect"
	"strings"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval/internal"
	"gitlab.com/distributed_lab/urlval/internal/betterreflect"
)

var errNotSupportedParameter = errors.New("query parameter is not supported for this endpoint")

// Decode is decodes provided url values to destination struct.
// Using Decode requires your request to follow JSON API spec -  If it's
// values contains anything except from "include", "sort", "search",  "filter", "page",
// or query has parameters that are not tagged in dest -  Decode still populates
// dest but also returns an error. The only error type it returns is errBadRequest
// that is (hopefully) compatible with ape  (https://gitlab.com/distributed_lab/ape),
// so can be rendered directly to client.
func Decode(values url.Values, dest interface{}) error {
	tokens, _ := internal.Tokenize(values)
	refdest := betterreflect.NewStruct(dest)
	setDefaults(refdest)
	errs := errBadRequest{}

	for token := range tokens {
		if token.Type == internal.TokenTypeInvalid {
			errs[token.Key] = errNotSupportedParameter
			continue
		}

		ok, err := decodeToken(token, refdest)
		if err != nil {
			if errors.Cause(err) == betterreflect.ErrInvalidType {
				errs[token.Key] = err
				continue
			}

			return errors.Wrap(err, "failed to decode token")
		}
		if !ok {
			errs[token.Key] = errNotSupportedParameter
		}
	}

	return errs.Filter()
}

func setDefaults(s *betterreflect.Struct) {
	for i := 0; i < s.NumField(); i++ {
		if s.Value(i).Kind() == reflect.Struct {
			nestedStruct := betterreflect.NewStructFromValue(s.Value(i))
			setDefaults(nestedStruct)
			continue
		}

		if betterreflect.IsZero(s.Value(i).Interface()) && s.Tag(i, "default") != "" {
			if err := s.Set(i, s.Tag(i, "default")); err != nil {
				panic(errors.Wrap(err, "failed to set default value"))
			}
		}
	}
}

func decodeToken(token internal.Token, dest *betterreflect.Struct) (bool, error) {
	var decoded bool

	for i := 0; i < dest.NumField(); i++ {
		var ok bool
		var err error

		if dest.Type(i).Kind() == reflect.Struct {
			if ok, err = decodeToken(token, betterreflect.NewStructFromValue(dest.Value(i))); err != nil {
				return false, err
			}
			if ok {
				if decoded {
					panic(errors.New("decoding same token twice - probably your struct has 2 or more similar tags"))
				}

				decoded = true
				// not returning here, because we still need to traverse
				// whole struct to ensure we don't decode same token twice.
			}

			continue
		}

		if ok, err = trySet(dest, i, token); err != nil {
			return false, err
		}

		if ok {
			if decoded {
				panic(errors.New("decoding same token twice - probably your struct has 2 or more similar tags"))
			}

			decoded = true
			// not returning here, because we still need to traverse
			// whole struct to ensure we don't decode same token twice.
		}
	}

	return decoded, nil
}

func trySet(dest *betterreflect.Struct, i int, token internal.Token) (bool, error) {
	var value interface{}

	switch token.Type {
	case internal.TokenTypeInclude:
		if dest.Tag(i, "include") == token.Key {
			if dest.Type(i).Kind() != reflect.Bool {
				panic("invalid destination type, expected bool for include tags")
			}

			value = true
		}
	case internal.TokenTypeFilter:
		if dest.Tag(i, "filter") == token.Key {
			if dest.Type(i).Kind() != reflect.Ptr &&
				dest.Type(i).Kind() != reflect.Slice &&
				dest.Type(i).Elem().Kind() != reflect.String {
				panic("invalid destination type, expected pointer or []string for filter tags")
			}

			value = token.Value
		}
	case internal.TokenTypePage:
		if dest.Tag(i, "page") == token.Key {
			value = token.Value
		}
	case internal.TokenTypeSearch:
		if dest.Tag(i, "url") == "search" {
			if dest.Type(i).Kind() != reflect.Ptr || dest.Type(i).Elem().Kind() != reflect.String {
				panic("invalid destination type, expected *string for search tags")
			}

			value = token.Value
		}
	case internal.TokenTypeSort:
		if dest.Tag(i, "url") == "sort" {
			if dest.Type(i).Kind() != reflect.Slice || dest.Type(i).Elem().Kind() != reflect.String {
				panic("invalid destination type, expected []string (or any alias with underlying []string) for sort tags")
			}
			value = strings.Split(token.Value, ",")
		}
	default:
		panic(errors.Errorf("unknown token type: %d", token.Type))
	}

	if value == nil {
		return false, nil
	}

	if err := dest.Set(i, value); err != nil {
		return false, err
	}

	return true, nil
}
