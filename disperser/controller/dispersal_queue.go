package controller

// Responsible for polling dynamo for blobs that need to be put into batches for dispersal.
type DispersalQueue struct {
}

func NewDispersalQueue() (*DispersalQueue, error) {
	return &DispersalQueue{}, nil
}

// Close the queue and free any resources.
func (q *DispersalQueue) Close() error {
	return nil
}
