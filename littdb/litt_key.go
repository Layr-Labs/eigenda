package littdb

// LittKey is a key in the LittDB.
type LittKey struct {
	// Table is the table in which the key resides. Two LitKeys with the same Key but different Tables
	// do not collide with each other.
	Table string
	// Key is the key in the table.
	Key []byte
}
