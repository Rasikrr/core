package util

import "reflect"

// IsNil возвращает true только для nil-значений nilable-типов:
// *T, []T, map[K]V, chan T, func, interface{}, unsafe.Pointer.
// Корректно обрабатывает кейс: var e error = (*MyErr)(nil).
func IsNil[T any](v T) bool {
	if any(v) == nil {
		return true
	}
	rv := reflect.ValueOf(any(v))

	// распаковываем интерфейсы до конкретного значения
	for rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return true
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan, reflect.UnsafePointer:
		return rv.IsNil()
	default:
		return false
	}
}

func IsZero[T any](v T) bool {
	return reflect.ValueOf(v).IsZero()
}
