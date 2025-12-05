// Package util contains helper utilities such as phone normalization.
package util

import "github.com/nyaruka/phonenumbers"

// NormalizeE164 attempts to format the number to E.164 using the provided region (e.g. "US").
// If parsing fails, it returns the original input.
func NormalizeE164(input, region string) string {
	num, err := phonenumbers.Parse(input, region)
	if err != nil || !phonenumbers.IsValidNumber(num) {
		return input
	}
	return phonenumbers.Format(num, phonenumbers.E164)
}
