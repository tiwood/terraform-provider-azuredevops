package converter

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf16"

	"github.com/microsoft/azure-devops-go-api/azuredevops/licensing"
)

// String Get a pointer to a string
func String(value string) *string {
	if strings.EqualFold(value, "") {
		return nil
	}
	return &value
}

// StringFromInterface get a string pointer from an interface
func StringFromInterface(value interface{}) *string {
	return String(value.(string))
}

// Bool Get a pointer to a boolean value
func Bool(value bool) *bool {
	return &value
}

// Int Get a pointer to an integer value
func Int(value int) *int {
	return &value
}

// UInt64 Get a pointer to an uint64 value
func UInt64(value uint64) *uint64 {
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

func DecodeUtf16HexString(message string) (string, error) {
	b, err := hex.DecodeString(message)
	if err != nil {
		return "", err
	}
	ints := make([]uint16, len(b)/2)
	if err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &ints); err != nil {
		return "", err
	}
	return string(utf16.Decode(ints)), nil
}

func EncodeUtf16HexString(message string) (string, error) {
	runeByte := []rune(message)
	encodedByte := utf16.Encode(runeByte)
	var sb strings.Builder
	for i := 0; i < len(encodedByte); i++ {
		fmt.Fprintf(&sb, "%02x%02x", encodedByte[i], encodedByte[i]>>8)
	}
	return sb.String(), nil
}
