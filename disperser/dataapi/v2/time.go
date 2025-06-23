package v2

import "time"

const timeFormat = time.RFC3339Nano

// Format a time in [timeFormat] format for use in query parameters.
// This ensures that the server can parse the time correctly.
// Used for blobs and batches queries.
func FormatQueryParamTime(time time.Time) string {
	// Note that we need to convert to UTC() such that it gets formatted to
	// something like "2023-10-01T12:34:56.789Z" instead of "2023-10-01T12:34:56.789+00:00",
	// because `+` gets converted to a space in query parameters,
	// which is then not parsable as a RFC3339Nano time.
	return time.UTC().Format(timeFormat)
}

// Parse the time string in RFC3339Nano [timeFormat] format.
// This is used for parsing query parameters like "before" and "after",
// for blobs and batches queries.
// Meant to parse query params that are formatted with [FormatQueryParamTime].
func parseQueryParamTime(timeStr string) (time.Time, error) {
	return time.Parse(timeFormat, timeStr)
}
