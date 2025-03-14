package segment

import (
	"fmt"
	"strings"
)

// TODO this is a place holder. Every time a segment is sealed, write this file.

// WriteInfoFile writes the [segment-index].info file for the segment. This is a human-readable yaml formatted file
// intended for debugging and inspection purposes. It is not read/utilized by the DB, and does not have strong
// consistency guarantees if there is a crash.
func WriteInfoFile(
	targetDirectory string,
	index uint32,
	shardingFactor uint32,
	salt uint32,
	firstDataTimestamp uint64,
	lastDataTimestamp uint64,
	keyCount uint64,
	keyFileSize uint64,
	dataFileSize uint64,
	infoTimestamp uint64,
) error {

	return nil
}

// generateInfoString generates the human-readable string that will be written to the X.info file.
func generateInfoString(
	index uint32,
	sealed bool,
	shardingFactor uint32,
	salt uint32,
	firstDataTimestamp uint64,
	lastDataTimestamp uint64,
	keyCount uint64,
	keyFileSize uint64,
	dataFileSize uint64,
	infoTimestamp uint64,
) string {

	sb := strings.Builder{}

	sb.WriteString("# This file contains human-readable information about a segment.\n")
	sb.WriteString("# Its purpose is to aid debugging and to give insight into the DB.\n")
	sb.WriteString("# This file does not have strong consistency guarantees in the")
	sb.WriteString("# advent of a crash.\n\n")

	sb.WriteString(fmt.Sprintf("index: %d\n", index))
	sb.WriteString(fmt.Sprintf("sealed: %t\n", sealed))
	sb.WriteString(fmt.Sprintf("shardingFactor: %d\n", shardingFactor))
	sb.WriteString(fmt.Sprintf("salt: %d\n", salt))
	sb.WriteString(fmt.Sprintf("firstDataTimestamp: %d\n", firstDataTimestamp))
	sb.WriteString(fmt.Sprintf("lastDataTimestamp: %d\n", lastDataTimestamp))
	sb.WriteString(fmt.Sprintf("keyCount: %d\n", keyCount))
	sb.WriteString(fmt.Sprintf("keyFileSize: %d\n", keyFileSize))
	sb.WriteString(fmt.Sprintf("dataFileSize: %d\n", dataFileSize))
	sb.WriteString(fmt.Sprintf("infoTimestamp: %d\n", infoTimestamp))

	return sb.String()
}
