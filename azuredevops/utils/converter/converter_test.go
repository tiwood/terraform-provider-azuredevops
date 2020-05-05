// +build all utils converter

package converter

import (
	"testing"
)

func TestString(t *testing.T) {
	value := "Hello World"
	valuePtr := String(value)
	if value != *valuePtr {
		t.Errorf("The pointer returned references a different value")
	}
}

func TestInt(t *testing.T) {
	value := 123456
	valuePtr := Int(value)
	if value != *valuePtr {
		t.Errorf("The pointer returned references a different value")
	}
}

func TestBoolTrue(t *testing.T) {
	value := true
	valuePtr := Bool(value)
	if value != *valuePtr {
		t.Errorf("The pointer returned references a different value")
	}
}

func TestBoolFalse(t *testing.T) {
	value := false
	valuePtr := Bool(value)
	if value != *valuePtr {
		t.Errorf("The pointer returned references a different value")
	}
}
<<<<<<< HEAD
=======

func TestLicenseTypeAccount(t *testing.T) {
	assertAccountLicenseType(t, licensing.AccountLicenseTypeValues.None)
	assertAccountLicenseType(t, licensing.AccountLicenseTypeValues.EarlyAdopter)
	assertAccountLicenseType(t, licensing.AccountLicenseTypeValues.Advanced)
	assertAccountLicenseType(t, licensing.AccountLicenseTypeValues.Professional)
	assertAccountLicenseType(t, licensing.AccountLicenseTypeValues.Express)
	assertAccountLicenseType(t, licensing.AccountLicenseTypeValues.Professional)

	_, err := AccountLicenseType("foo")
	assert.Equal(t, err.Error(), "Error unable to match given AccountLicenseType:foo")
}

func assertAccountLicenseType(t *testing.T, accountLicenseType licensing.AccountLicenseType) {
	actualAccountLicenseType, err := AccountLicenseType(string(accountLicenseType))
	assert.Nil(t, err, fmt.Sprintf("Error should not thrown by %s", string(accountLicenseType)))
	assert.Equal(t, &accountLicenseType, actualAccountLicenseType, fmt.Sprintf("%s should be able to convert into the AccountLicenseType", string(accountLicenseType)))
}

func TestStringFromInterface_StringValue(t *testing.T) {
	value := "Hello World"
	valuePtr := StringFromInterface(value)
	if value != *valuePtr {
		t.Errorf("The pointer returned references a different value")
	}
}

func TestStringFromInterface_InterfaceValue(t *testing.T) {
	value := "Hello World"
	var interfaceValue interface{}

	interfaceValue = value
	valuePtr := StringFromInterface(interfaceValue)
	if value != *valuePtr {
		t.Errorf("The pointer returned references a different value")
	}
}

type encodeTestType struct {
	plainString   string
	encodedString string
}

var encodeTestCases = []encodeTestType{
	{
		plainString:   "branch_1_1",
		encodedString: "6200720061006e00630068005f0031005f003100",
	},
	{
		plainString:   "master",
		encodedString: "6d0061007300740065007200",
	},
}

func TestDecodeUtf16HexString(t *testing.T) {
	for _, etest := range encodeTestCases {
		val, err := DecodeUtf16HexString(etest.encodedString)
		assert.Nil(t, err, fmt.Sprintf("Error should not thrown by %s", etest.encodedString))
		assert.EqualValues(t, etest.plainString, val)
	}
}

func TestEncodeUtf16HexString(t *testing.T) {
	for _, etest := range encodeTestCases {
		val, err := EncodeUtf16HexString(etest.plainString)
		assert.Nil(t, err, fmt.Sprintf("Error should not thrown by %s", etest.plainString))
		assert.EqualValues(t, etest.encodedString, val)
	}
}
>>>>>>> origin/r_permissions
