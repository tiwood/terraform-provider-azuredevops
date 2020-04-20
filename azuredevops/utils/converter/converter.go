package converter

import (
	"fmt"
	"strings"

	"github.com/microsoft/azure-devops-go-api/azuredevops/licensing"
)

// String Get a pointer to a string
func String(value string) *string {
	if strings.EqualFold(value, "") {
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
	switch strings.ToLower(accountLicenseTypeValue) {
	case "none":
		accountLicenseType = licensing.AccountLicenseTypeValues.None
	case "earlyadopter":
		accountLicenseType = licensing.AccountLicenseTypeValues.EarlyAdopter
	case "basic":
		fallthrough
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

// AccountLicensingSource convert a string value to a licensing.AccountLicenseType pointer
func AccountLicensingSource(licensingSourceValue string) (*licensing.LicensingSource, error) {
	var licensingSource licensing.LicensingSource
	switch strings.ToLower(licensingSourceValue) {
	case "none":
		licensingSource = licensing.LicensingSourceValues.None
	case "account":
		licensingSource = licensing.LicensingSourceValues.Account
	case "msdn":
		licensingSource = licensing.LicensingSourceValues.Msdn
	case "profile":
		licensingSource = licensing.LicensingSourceValues.Profile
	case "auto":
		licensingSource = licensing.LicensingSourceValues.Auto
	case "trial":
		licensingSource = licensing.LicensingSourceValues.Trial
	default:
		return nil, fmt.Errorf("Error unable to match given LicensingSource :%s", licensingSourceValue)
	}
	return &licensingSource, nil
}
