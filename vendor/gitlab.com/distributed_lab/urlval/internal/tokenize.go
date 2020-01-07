package internal

import (
	"net/url"
	"regexp"
	"strings"
)

type TokenType int32

const (
	TokenTypeFilter TokenType = 1 + iota
	TokenTypeInclude
	TokenTypePage
	TokenTypeSort
	TokenTypeSearch
	TokenTypeInvalid
)

type Token struct {
	Type  TokenType
	Key   string
	Value string
}

// bool for consumed
type Tokens map[Token]bool

func Tokenize(values url.Values) (Tokens, error) {
	tokens := Tokens{}

	tokenizeIncludes(values, tokens)
	tokenizeFilters(values, tokens)
	tokenizePagination(values, tokens)
	tokenizeSort(values, tokens)
	tokenizeSearch(values, tokens)

	for k := range values {
		tokens[Token{
			Type: TokenTypeInvalid,
			Key:  k,
		}] = false
	}

	return tokens, nil
}

func tokenizeIncludes(values url.Values, tokens Tokens) {
	includes := values.Get("include")
	if includes == "" {
		return
	}
	for _, include := range strings.Split(includes, ",") {
		if include != "" { // so include=a,b,c, (with comma at the end) wont produce a 4th token
			tokens[Token{
				Type: TokenTypeInclude,
				Key:  include,
			}] = false
		}
	}
	values.Del("include")
}

func tokenizeFilters(values url.Values, tokens Tokens) {
	for k, v := range values {
		// TODO: check v length
		ok, key := extractFilter(k)
		if ok {
			tokens[Token{
				Type:  TokenTypeFilter,
				Key:   key,
				Value: v[0],
			}] = false
			values.Del(k)
		}
	}
}

func tokenizeSort(values url.Values, tokens Tokens) {
	sorts := values.Get("sort")
	if sorts == "" {
		values.Del("sort")
		return
	}

	tokens[Token{
		Type:  TokenTypeSort,
		Key:   "sort",
		Value: sorts,
	}] = false

	values.Del("sort")
}

func tokenizeSearch(values url.Values, tokens Tokens) {
	search := values.Get("search")
	if search == "" {
		return
	}

	tokens[Token{
		Type:  TokenTypeSearch,
		Key:   "search",
		Value: search,
	}] = false

	values.Del("search")
}

func extractFilter(s string) (bool, string) {
	r := regexp.MustCompile(`^filter\[([^\]]+)\]$`)
	match := r.FindStringSubmatch(s)
	if len(match) != 2 {
		return false, ""
	}
	return true, match[1]
}

func tokenizePagination(values url.Values, tokens Tokens) {
	for k, v := range values {
		// TODO: check v length
		ok, key := extractPage(k)
		if ok {
			tokens[Token{
				Type:  TokenTypePage,
				Key:   key,
				Value: v[0],
			}] = false
			values.Del(k)
		}
	}
}

func extractPage(s string) (bool, string) {
	r := regexp.MustCompile(`^page\[([^\]]+)\]$`)
	match := r.FindStringSubmatch(s)
	if len(match) != 2 {
		return false, ""
	}
	return true, match[1]
}
