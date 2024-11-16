package files

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"

	"google.golang.org/genproto/googleapis/api/httpbody"
)

func ParseMultipartForm[T any](ctx context.Context, body *httpbody.HttpBody) (*T, error) {
	dst := new(T)

	formData, err := NewFormData(ctx, body)
	if err != nil {
		return dst, err
	}

	data := make(map[string][]any)

	for key, values := range formData.form.Value {
		if len(values) > 0 {
			for _, v := range values {
				data[key] = append(data[key], v)
			}
		}
	}

	for key, files := range formData.form.File {
		if len(files) > 0 {
			for _, fileHeader := range files {
				data[key] = append(data[key], fileHeader)
			}
		}
	}

	err = mapToStruct(data, dst)
	if err != nil {
		return dst, err
	}

	return dst, nil
}

func mapToStruct(data map[string][]any, dst any) (err error) {
	dstVal := reflect.ValueOf(dst).Elem()

	if dstVal.Kind() == reflect.Struct {
		for i := 0; i < dstVal.NumField(); i++ {
			field := dstVal.Type().Field(i)
			fieldVal := dstVal.Field(i)
			key := field.Tag.Get("form")

			if key == "" {
				key = field.Name
			}

			if value, exists := data[key]; exists {
				var valueType reflect.Type
				if len(value) > 0 {
					valueType = reflect.TypeOf(value[0])
				}

				if field.Type.Kind() == reflect.Ptr || field.Type.Kind() == reflect.Pointer { //nolint:gocritic
					if field.Type.Elem().Kind() != reflect.Slice && field.Type.Elem().Kind() != reflect.Array {
						err = setFieldValue(fieldVal, value[0])
					}
				} else if field.Type.Kind() != reflect.Slice && field.Type.Kind() != reflect.Array {
					err = setFieldValue(fieldVal, value[0])
				} else {
					err = setFieldValue(fieldVal, value, valueType)
				}

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func setFieldValue(fieldVal reflect.Value, valueStr any, dataValType ...reflect.Type) error {
	switch fieldVal.Kind() {
	case reflect.Ptr:
		if fieldVal.IsNil() {
			fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
		}

		return setFieldValue(fieldVal.Elem(), valueStr)

	case reflect.Struct:
		if str, ok := valueStr.(string); ok {
			return json.Unmarshal([]byte(str), fieldVal.Addr().Interface())
		}

		val := reflect.ValueOf(valueStr)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		fieldVal.Set(val)

	case reflect.String:
		fieldVal.SetString(valueStr.(string))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setIntValue(fieldVal, valueStr.(string))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setUintValue(fieldVal, valueStr.(string))

	case reflect.Float32, reflect.Float64:
		return setFloatValue(fieldVal, valueStr.(string))

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(valueStr.(string))
		if err != nil {
			return errors.New("failed to parse bool")
		}

		fieldVal.SetBool(boolValue)

	case reflect.Slice, reflect.Array:
		if fieldVal.Type().Elem().Kind() == reflect.Uint8 {
			fieldVal.SetBytes([]byte(valueStr.(string)))
		} else {
			// create a new slice of the same type as the data
			if len(dataValType) == 0 {
				return errors.New("Unsupported slice type " + fieldVal.Type().Elem().Kind().String())
			}

			sliceType := reflect.SliceOf(dataValType[0])
			newSlice := reflect.MakeSlice(sliceType, 0, 0)

			for _, v := range valueStr.([]any) {
				newSlice = reflect.Append(newSlice, reflect.ValueOf(v))
			}

			fieldVal.Set(newSlice)
		}
	case reflect.Interface:
		fieldVal.Set(reflect.ValueOf(valueStr))
	case reflect.Map:
		if fieldVal.Type().Key().Kind() == reflect.String {
			if err := json.Unmarshal([]byte(valueStr.(string)), fieldVal.Addr().Interface()); err != nil {
				return errors.New("failed to unmarshal map")
			}
		} else {
			return errors.New("Unsupported map key type " + fieldVal.Type().Key().Kind().String())
		}

	case reflect.Invalid, reflect.Uintptr, reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func,
		reflect.UnsafePointer:
		fallthrough
	default:
		return errors.New("Unsupported kind " + fieldVal.Kind().String())
	}

	return nil
}

func setIntValue(fieldVal reflect.Value, valueStr string) error {
	intValue, err := strconv.ParseInt(valueStr, 10, fieldVal.Type().Bits())
	if err != nil {
		return errors.New("failed to parse int")
	}

	fieldVal.SetInt(intValue)

	return nil
}

func setUintValue(fieldVal reflect.Value, valueStr string) error {
	uintValue, err := strconv.ParseUint(valueStr, 10, fieldVal.Type().Bits())
	if err != nil {
		return errors.New("failed to parse uint")
	}

	fieldVal.SetUint(uintValue)

	return nil
}

func setFloatValue(fieldVal reflect.Value, valueStr string) error {
	floatValue, err := strconv.ParseFloat(valueStr, fieldVal.Type().Bits())
	if err != nil {
		return errors.New("failed to parse float")
	}

	fieldVal.SetFloat(floatValue)

	return nil
}
