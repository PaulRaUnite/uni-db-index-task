package urlval

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval/internal/betterreflect"
)

func Encode(src interface{}) string {
	refsrc := betterreflect.NewStruct(src)
	values := url.Values{}

	populateValues(values, refsrc)

	return values.Encode()
}

func populateValues(values url.Values, refsrc *betterreflect.Struct) {
	for i := 0; i < refsrc.NumField(); i++ {
		if refsrc.Type(i).Kind() == reflect.Struct {
			populateValues(values, betterreflect.NewStructFromValue(refsrc.Value(i)))
			continue
		}

		switch name, tag := getTag(refsrc, i); name {
		case "include":
			if refsrc.Value(i).Bool() {
				includes := values.Get("include")
				if includes == "" {
					includes = tag
				} else {
					includes += "," + tag
				}
				setValue(values, "include", includes)
			}
		case "page":
			key := fmt.Sprintf("page[%s]", tag)
			value := refsrc.Value(i)
			setValue(values, key, value)
		case "filter":
			key := fmt.Sprintf("filter[%s]", tag)
			setValue(values, key, refsrc.Value(i))
		case "url":
			switch tag {
			case "sort":
				if sorts, err := betterreflect.ConvertSlice(refsrc.Value(i).Interface(), reflect.TypeOf("")); err == nil {
					value := strings.Join(sorts.([]string), ",")
					setValue(values, "sort", value)
				}
			case "search":
				value := toString(refsrc.Value(i).Interface())
				setValue(values, "search", value)
			}
		}

	}
}

func setValue(values url.Values, key string, value interface{}) {
	if strVal := toString(value); strVal != "" {
		values.Set(key, strVal)
	}
}

func toString(value interface{}) string {
	if v, ok := value.(reflect.Value); ok {
		value = v.Interface()
	}

	if v, ok := value.([]string); ok {
		value = strings.Join(v, ",")
	}

	// some magic to convert values of custom aliased types to their,
	// underlying type, because cast fails to do this:
	if value = betterreflect.ConvertToUnderlyingType(value); value == nil {
		return ""
	}

	return cast.ToString(value)
}

func getTag(refsrc *betterreflect.Struct, i int) (name, tag string) {
	var names = []string{
		"include",
		"page",
		"filter",
		"sort",
		"url",
	}
	for _, name = range names {
		if tag = refsrc.Tag(i, name); tag != "" {
			return name, tag
		}
	}

	return "", ""
}
