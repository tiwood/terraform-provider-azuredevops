// +build all helper converter

package converter

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/microsoft/azure-devops-go-api/azuredevops/licensing"
	"github.com/stretchr/testify/assert"
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

type dataStruct struct {
	Name        string
	Value       string
	ValueRef    reflect.Value
	ValuePtr    *string
	ValueRefPtr *reflect.Value
}

var value1 = "Value1Ptr"
var value1Ref = reflect.ValueOf("Value1RefPtr")
var value2 = "Value2Ptr"
var value2Ref = reflect.ValueOf("Value2RefPtr")
var value3 = "Value3Ptr"
var value3Ref = reflect.ValueOf("Value3RefPtr")

var convDataStructTest = []dataStruct{
	dataStruct{
		Name:        "Name1",
		Value:       "Value1",
		ValueRef:    reflect.ValueOf("Value1Ref"),
		ValuePtr:    &value1,
		ValueRefPtr: &value1Ref,
	},
	dataStruct{
		Name:        "Name2",
		Value:       "Value2",
		ValueRef:    reflect.ValueOf("Value2Ref"),
		ValuePtr:    &value2,
		ValueRefPtr: &value2Ref,
	},
	dataStruct{
		Name:        "Name3",
		Value:       "Value3",
		ValueRef:    reflect.ValueOf("Value3Ref"),
		ValuePtr:    &value3,
		ValueRefPtr: &value3Ref,
	},
	dataStruct{
		Name:        "Name4",
		Value:       "Value1",
		ValueRef:    reflect.ValueOf("Value1Ref"),
		ValuePtr:    &value1,
		ValueRefPtr: &value1Ref,
	},
	dataStruct{
		Name:        "Name5",
		Value:       "Value1",
		ValueRef:    reflect.ValueOf("Value1Ref"),
		ValuePtr:    &value1,
		ValueRefPtr: &value1Ref,
	},
}

func TestGetValueByName(t *testing.T) {

}

func TestGetValueSliceByName(t *testing.T) {
}

func TestFilterObjectsByAttributeValues_Struct_DirectValue(t *testing.T) {
	attrComp := []AttributeComparison{
		AttributeComparison{
			Name:  "Value",
			Value: "Value1",
		},
	}
	ret, err := FilterObjectsByAttributeValues(convDataStructTest, &attrComp)
	assert.Nil(t, err)
	arr, ok := ret.([]dataStruct)
	assert.True(t, ok)
	assert.Len(t, arr, 3)
}

func TestFilterObjectsByAttributeValues_Struct_PtrValue(t *testing.T) {
	attrComp := []AttributeComparison{
		AttributeComparison{
			Name:  "ValuePtr",
			Value: "Value1Ptr",
		},
	}
	ret, err := FilterObjectsByAttributeValues(convDataStructTest, &attrComp)
	assert.Nil(t, err)
	arr, ok := ret.([]dataStruct)
	assert.True(t, ok)
	assert.Len(t, arr, 3)
}

func TestFilterObjectsByAttributeValues_Struct_RefValue(t *testing.T) {
	attrComp := []AttributeComparison{
		AttributeComparison{
			Name:  "ValueRef",
			Value: "Value1Ref",
		},
	}
	ret, err := FilterObjectsByAttributeValues(convDataStructTest, &attrComp)
	assert.Nil(t, err)
	arr, ok := ret.([]dataStruct)
	assert.True(t, ok)
	assert.Len(t, arr, 3)
}

func TestFilterObjectsByAttributeValues_Struct_RefPtrValue(t *testing.T) {
	attrComp := []AttributeComparison{
		AttributeComparison{
			Name:  "ValueRefPtr",
			Value: "Value1RefPtr",
		},
	}
	ret, err := FilterObjectsByAttributeValues(convDataStructTest, &attrComp)
	assert.Nil(t, err)
	arr, ok := ret.([]dataStruct)
	assert.True(t, ok)
	assert.Len(t, arr, 3)
}
