// Package pressure provides the Pressure data-type and helper functions
package pressure

import "fmt"

// Pressure data type
type Pressure struct {
	Systolic  int   `json:"systolic"`
	Diastolic int   `json:"diastolic"`
	HeartRate int   `json:"heartrate,omitempty"`
	Timestamp int64 `json:"timestamp"`
}

func (p Pressure) String() string {
	return fmt.Sprintf("{systolic:%d, diastolic:%d, heartrate:%d, timestamp:%d}", p.Systolic, p.Diastolic, p.HeartRate, p.Timestamp)
}

// Valid checks if a Pressure is valid
func (p Pressure) Valid() bool {
	return p.Timestamp > 0 && p.Diastolic > 0 && p.Systolic > 0
}
