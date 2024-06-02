package article

import "time"

// CreationDate is the time in the past when the Article was created.
type CreationDate time.Time

// HumanReadableString converts CreationDate to a human-readable string.
// The date format is RFC822 chosen as the most common date format.
func (cd CreationDate) HumanReadableString() string {
	t := time.Time(cd)
	localTime := t.Local()
	return localTime.Format(time.RFC822)
}
