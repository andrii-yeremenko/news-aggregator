package article

import "time"

// CreationDate is the time in the past when the Article was created.
type CreationDate time.Time

// HumanReadableString converts CreationDate to a human-readable string.
// The date format is RFC822 chosen as a most common date format.
func (cd CreationDate) HumanReadableString() string {
	t := time.Time(cd)
	return t.Format(time.RFC822)
}
