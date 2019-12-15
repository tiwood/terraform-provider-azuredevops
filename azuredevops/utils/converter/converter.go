package converter

import (
	"crypto/sha1"
	"fmt"
	"reflect"
	"strings"

	"github.com/microsoft/azure-devops-go-api/azuredevops/licensing"
)

// String Get a pointer to a string
func String(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// Bool Get a pointer to a boolean value
func Bool(value bool) *bool {
	return &value
}

// Int Get a pointer to an integer value
func Int(value int) *int {
	return &value
}

// ToString Given a pointer return its value, or a default value of the poitner is nil
func ToString(value *string, defaultValue string) string {
	if value != nil {
		return *value
	}

	return defaultValue
}

// ToBool Given a pointer return its value, or a default value of the pointer is nil
func ToBool(value *bool, defaultValue bool) bool {
	if value != nil {
		return *value
	}

	return defaultValue
}

// AccountLicenseType Get a pointer to an AccountLicenseType
func AccountLicenseType(accountLicenseTypeValue string) (*licensing.AccountLicenseType, error) {
	var accountLicenseType licensing.AccountLicenseType
	switch accountLicenseTypeValue {
	case "none":
		accountLicenseType = licensing.AccountLicenseTypeValues.None
	case "earlyAdopter":
		accountLicenseType = licensing.AccountLicenseTypeValues.EarlyAdopter
	case "express":
		accountLicenseType = licensing.AccountLicenseTypeValues.Express
	case "professional":
		accountLicenseType = licensing.AccountLicenseTypeValues.Professional
	case "advanced":
		accountLicenseType = licensing.AccountLicenseTypeValues.Advanced
	case "stakeholder":
		accountLicenseType = licensing.AccountLicenseTypeValues.Stakeholder
	default:
		return nil, fmt.Errorf("Error unable to match given AccountLicenseType:%s", accountLicenseTypeValue)
	}
	return &accountLicenseType, nil
}

// ConvertToStringSlice convert a slice to a slice containing strings
func ToStringSlice(input []interface{}) []string {
	result := make([]string, len(input))
	for i, k := range input {
		result[i] = k.(string)
	}

	return result
}

// GetValue returns a direct value
func GetValue(input interface{}) interface{} {
	var s reflect.Value
	if reflect.TypeOf(input) == reflect.TypeOf((*reflect.Value)(nil)).Elem() {
		s = input.(reflect.Value)
	} else {
		s = reflect.ValueOf(input)
	}
	if s.Kind() == reflect.Ptr {
		if s.IsNil() {
			return nil
		}
		s = s.Elem()
	}
	if s.Kind() == reflect.Ptr && s.IsNil() {
		return nil
	}
	return s.Interface()
}

// GetValueByName returns a value from a structured type like mape and struct
func GetValueByName(input interface{}, name string) interface{} {
	/*
		var s reflect.Value
		if reflect.TypeOf(input) == reflect.TypeOf((*reflect.Value)(nil)).Elem() {
			s = input.(reflect.Value)
		} else {
			s = reflect.ValueOf(input)
		}
		if s.Kind() == reflect.Ptr {
			if s.IsNil() {
				return nil
			}
			s = s.Elem()
		}
	*/
	v := GetValue(input)
	s := reflect.ValueOf(v)

	if s.Kind() == reflect.Struct {
		f := s.FieldByName(name)
		if !f.IsValid() {
			panic(fmt.Sprintf("Struct does not contain property %s", name))
		}
		if f.Kind() == reflect.Ptr {
			if f.IsNil() {
				return nil
			}
			f = reflect.Indirect(f)
		}
		return f.Interface()
	} else if s.Kind() == reflect.Map {
		ifc := s.Interface()
		if imap, ok := ifc.(map[string]interface{}); ok {
			return imap[name]
		}
		panic(fmt.Sprintf("Map %T must be of form map[string]interface{}", s))
	}
	panic(fmt.Sprintf("Type %T is not a structured type (struct, map)", s))
}

// GetValueSliceByName returns a slice of values intified by an attribute name
func GetValueSliceByName(input *[]interface{}, attributeName string) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, len(*input))
	for i, user := range *input {
		output[i] = GetValueByName(user, attributeName)
	}
	return output
}

// AttributeComparison defines a comparison on an (struct) attribute
type AttributeComparison struct {
	Name       string
	Value      string
	IgnoreCase bool
	AllowNil   bool
}

// FilterObjectsByAttributeValues returns a filtered slice of objects by an array of comparisons
func FilterObjectsByAttributeValues(input interface{}, comparison *[]AttributeComparison) (interface{}, error) {
	if comparison == nil || len(*comparison) <= 0 {
		t := reflect.TypeOf(input)
		if t.Kind() == reflect.Ptr {
			return reflect.Indirect(reflect.ValueOf(input)).Interface(), nil
		}
		return input, nil
	}

	vi := reflect.ValueOf(input)
	t := reflect.TypeOf(input)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Slice {
		panic("Input must be a slice")
	}
	vt := t.Elem()
	output := reflect.MakeSlice(reflect.SliceOf(vt), 0, 0)
	if reflect.TypeOf(input).Kind() != reflect.Ptr || !vi.IsNil() {
		s := reflect.Indirect(vi)
		for i := 0; i < s.Len(); i++ {
			user := s.Index(i)
			b := true
			for _, comp := range *comparison {
				v := GetValue(GetValueByName(user, comp.Name))
				if v == nil {
					if comp.AllowNil {
						continue
					} else {
						b = false
						break
					}
				}
				if comp.IgnoreCase {
					if !strings.EqualFold(comp.Value, v.(string)) {
						b = false
						break
					}
				} else {
					if comp.Value != v.(string) {
						b = false
						break
					}
				}
			}
			if b {
				output = reflect.Append(output, user)
			}
		}
	}
	return output.Interface(), nil
}

// ToSHA1Hash returns a SHA1 hash code of a string slice, where the elements are concatenated using '-' as separator at first
func ToSHA1Hash(input *[]string) ([]byte, error) {
	h := sha1.New()
	if input != nil {
		if _, err := h.Write([]byte(strings.Join(*input, "-"))); err != nil {
			return nil, fmt.Errorf("Unable to compute hash for user descriptors: %v", err)
		}
	}
	return h.Sum(nil), nil
}
