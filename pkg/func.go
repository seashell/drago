package util

import "time"

// boolToPtr returns the pointer to a boolean
func BoolToPtr(b bool) *bool {
	return &b
}

// IntToPtr returns the pointer to an int
func IntToPtr(i int) *int {
	return &i
}

// Int8ToPtr returns the pointer to an int8
func Int8ToPtr(i int8) *int8 {
	return &i
}

// Int64ToPtr returns the pointer to an int
func Int64ToPtr(i int64) *int64 {
	return &i
}

// Uint64ToPtr returns the pointer to an uint64
func Uint64ToPtr(u uint64) *uint64 {
	return &u
}

// UintToPtr returns the pointer to an uint
func UintToPtr(u uint) *uint {
	return &u
}

// StringToPtr returns the pointer to a string
func StringToPtr(str string) *string {
	return &str
}

// TimeToPtr returns the pointer to a time stamp
func TimeToPtr(t time.Duration) *time.Duration {
	return &t
}

// Float64ToPtr returns the pointer to an float64
func Float64ToPtr(f float64) *float64 {
	return &f
}
