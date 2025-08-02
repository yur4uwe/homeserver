package test_utils

import "reflect"

// MockFunction replaces a global function with a mock function and returns a cleanup function.
func MockFunction(originalFuncPtr interface{}, mockFunc interface{}, cleanup func()) func() {
	// Validate that originalFuncPtr is a pointer to a function
	originalValue := reflect.ValueOf(originalFuncPtr)
	if originalValue.Kind() != reflect.Ptr || originalValue.Elem().Kind() != reflect.Func {
		panic("originalFuncPtr must be a pointer to a function")
	}

	// Save the original function
	originalFunc := originalValue.Elem()

	// Replace the original function with the mock function
	mockValue := reflect.ValueOf(mockFunc)
	originalValue.Elem().Set(mockValue)

	// Return a cleanup function to restore the original function
	return func() {
		originalValue.Elem().Set(originalFunc)
		cleanup()
	}
}

func MockGlobalVar(globalVarPtr interface{}, mockValue interface{}, cleanup func()) func() {
	// Validate that globalVarPtr is a pointer to a variable
	globalValue := reflect.ValueOf(globalVarPtr)
	if globalValue.Kind() != reflect.Ptr {
		panic("globalVarPtr must be a pointer to a variable")
	}

	// Save the original value
	originalValue := globalValue.Elem()

	// Set the mock value
	mockValueReflect := reflect.ValueOf(mockValue)
	globalValue.Elem().Set(mockValueReflect)

	// Return a cleanup function to restore the original value
	return func() {
		globalValue.Elem().Set(originalValue)
		cleanup()
	}
}
