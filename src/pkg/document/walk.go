package document

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type foundFunc func(v interface{}) interface{}

func setValue(v interface{}) foundFunc {
	return func(n interface{}) interface{} {
		// TODO don't panic
		// TODO only allow values to be set to the same type
		// if !reflect.ValueOf(n).CanSet() {
		// 	return n
		// }
		return v
	}
}

// walk walks along the root interface until the end of the path.
// It returns the found value or an error if the path is invalid.
func walk(path []string, root interface{}, f foundFunc) (interface{}, error) {
	current := root
	parent := root
	for pIdx, p := range path {
		parent = current
		v := reflect.ValueOf(current)

		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

	SwitchKind:
		switch v.Kind() {
		case reflect.Bool, reflect.Int, reflect.Float32, reflect.String:
			current = v.Interface()

		case reflect.Map:
			iter := v.MapRange()
			for iter.Next() {
				key := iter.Key()
				if key.Type().Kind() == reflect.String && key.String() != p {
					continue
				}
				current = iter.Value().Interface()
				if pIdx == len(path)-1 && f != nil {
					reflect.ValueOf(parent).SetMapIndex(key, reflect.ValueOf(f(current)))
				}
				break SwitchKind
			}
			return nil, fmt.Errorf("invalid path %v %+v", p, v)

		case reflect.Slice:
			sIdx, err := strconv.Atoi(p)
			if err != nil {
				// TODO better error
				return nil, fmt.Errorf("error walking array. index %v: %w", sIdx, err)
			}
			if sIdx > v.Len() || sIdx < 0 {
				return nil, fmt.Errorf("index %v out of bounds", sIdx)
			}
			current = v.Index(sIdx).Interface()
			if pIdx == len(path)-1 && f != nil {
				reflect.ValueOf(parent).Index(sIdx).Set(reflect.ValueOf(f(current)))
			}

		case reflect.Struct:
			field := v.FieldByNameFunc(func(f string) bool {
				return strings.ToLower(f) == strings.ToLower(p)
			})
			if field.Kind() == reflect.Invalid || field.IsZero() {
				return nil, fmt.Errorf("invalid path %q at %q in struct %+v", path, p, root)
			}
			current = field.Interface()
			if pIdx == len(path)-1 && f != nil {
				field.Set(reflect.ValueOf(f(current)))
			}

		default:
			return nil, fmt.Errorf("unknown type %T", current)
		}
	}

	return current, nil
}
