package indexer

type AccumulatorObject interface {
}

type Accumulator interface {
	InitializeObject(header Header) (AccumulatorObject, error)

	UpdateObject(object AccumulatorObject, header *Header, event Event) (AccumulatorObject, error)

	// Serialize object takes the accummulator object, and serializes it using the rules for the specified fork.
	SerializeObject(object AccumulatorObject, fork UpgradeFork) ([]byte, error)

	// DeSerialize object deserializes an accumulator object using the rules for the specified fork.
	DeserializeObject(data []byte, fork UpgradeFork) (AccumulatorObject, error)
}
