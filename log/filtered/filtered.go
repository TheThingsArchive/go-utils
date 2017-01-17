package filtered

import (
	"reflect"
	"strings"

	"github.com/TheThingsNetwork/go-utils/log"
)

// A Filter is something that can filter a field value
type Filter interface {
	Filter(string, interface{}) interface{}
}

// A FilterFunc is a function that implements the Filter interface
type FilterFunc func(string, interface{}) interface{}

// Filter implements the Filter interface for FilterFuncs
func (fn FilterFunc) Filter(k string, v interface{}) interface{} {
	return fn(k, v)
}

// Filtered is a logger that filters the fields it receives
type Filtered struct {
	log.Interface
	filters []Filter
}

// Wrap wraps an existing logger, filtering the fields as it goes along
// using the provided filters
func Wrap(logger log.Interface, filters ...Filter) *Filtered {
	return &Filtered{
		Interface: logger,
		filters:   filters,
	}
}

// WithFilter creates a new Filtered that will use the extra filters
func (f *Filtered) WithFilters(filters ...Filter) *Filtered {
	return &Filtered{
		Interface: f.Interface,
		filters:   append(f.filters, filters...),
	}
}

// WithField filters the field and passes it on to the wrapped loggers WithField
func (f *Filtered) WithField(k string, v interface{}) log.Interface {
	val := v

	// apply the filters
	for _, filter := range f.filters {
		val = filter.Filter(k, val)
	}

	f.Interface = f.Interface.WithField(k, val)
	return f
}

// WithFields filters the fields and passes them on to the wrapped loggers WithFields
func (f *Filtered) WithFields(fields log.Fields) log.Interface {
	res := make(map[string]interface{}, len(fields))

	for k, v := range fields {
		val := v

		// apply the filters
		for _, filter := range f.filters {
			val = filter.Filter(k, val)
		}

		res[k] = val
	}

	f.Interface = f.Interface.WithFields(res)
	return f
}

var (
	defaultElided    = "<elided>"
	defaultSensitive = []string{
		"token",
		"access_token",
		"refresh_token",
		"key",
		"password",
		"code",
	}

	// DefaultSensitiveFilter is a Filter that filters most sensitive data and
	// replaces it with `<elided>`. These fields are filtered:
	// token, access_token, refresh_token, key, password, code
	DefaultSensitiveFilter = FilterSensitive(defaultSensitive, defaultElided)

	// QueryFilter filters maps at the fields with the name 'query', using the same rules as
	// DefaultSensitiveFilter
	DefaultQueryFilter = RestrictFilter("query", LowerCaseFilter(MapFilter(SliceFilter(DefaultSensitiveFilter))))
)

// FilterSensitive creates a Filter that filters most sensitive data like passwords,
// keys, access_tokens, etc. and replaces them with the elided value
func FilterSensitive(sensitive []string, elided interface{}) Filter {
	return FilterFunc(func(key string, v interface{}) interface{} {
		lower := strings.ToLower(key)
		for _, s := range sensitive {
			if lower == s {
				return elided
			}
		}

		return v
	})
}

// SliceFilter lifts the filter to also work on slices. It loses the
// type information of the slice elements
func SliceFilter(filter Filter) Filter {
	return FilterFunc(func(k string, v interface{}) interface{} {
		r := reflect.ValueOf(v)
		if r.Kind() == reflect.Slice {
			res := make([]interface{}, 0, r.Len())
			for i := 0; i < r.Len(); i++ {
				el := r.Index(i).Interface()
				res = append(res, filter.Filter(k, el))
			}

			return res
		}

		return filter.Filter(k, v)
	})
}

// MapFilter lifts the filter to also work on maps. It loses the type
// information of the map fields
func MapFilter(filter Filter) Filter {
	return FilterFunc(func(k string, v interface{}) interface{} {
		r := reflect.ValueOf(v)
		if r.Kind() == reflect.Map {
			// res will be the filtered map
			res := make(map[string]interface{}, r.Len())
			for _, key := range r.MapKeys() {
				str := key.String()
				val := r.MapIndex(key).Interface()
				res[str] = filter.Filter(str, val)
			}

			return res
		}

		return v
	})
}

// RestrictFilter restricts the filter to only work on a certain field
func RestrictFilter(fieldName string, filter Filter) Filter {
	return FilterFunc(func(k string, v interface{}) interface{} {
		if fieldName == k {
			return filter.Filter(k, v)
		}

		return v
	})
}

// LowerCaseFilter creates a filter that only get passed lowercase field names
func LowerCaseFilter(filter Filter) Filter {
	return FilterFunc(func(k string, v interface{}) interface{} {
		return filter.Filter(strings.ToLower(k), v)
	})
}
