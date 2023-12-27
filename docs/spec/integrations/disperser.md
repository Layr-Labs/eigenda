# Disperser

## Requests to Store Data

Requesters that want to store data on EigenDA (mostly rollups) make requests to the disperser with the form of [`BlobStoreRequest`](./types.md#blobstorerequest).

They specify the data they want to store on EigenDA and the different assumptions they are making on the quorums that they want to store their data. Requesters also attach a unique `BlobID` used to identify their blob later in the protocol. This is randomly generated and the disperser will fail the request if it has seen the given `BlobID` before. The disperser takes each `BlobStoreRequest` and adds it to a queue.

## Dispersal

### Dequeueing Requests

Every 30 seconds (configurable), the disperser takes all of the `BlobStoreRequest`s out of its queue, converts each `BlobStoreRequest`'s data into a polynomial [TODO: links @gpsanant @mooselumph], and "concatenates" the polynomials with each other using degree shifts.

In pseudocode:
```go
func DataToPoly(data []byte) ([]fr.Element) {
    poly := make([]fr.Element, 0)
    for i := 0; i <= len(data)/31; i++ {
        poly = append(poly, new(fr.Element).FromBytes(data[31*i:min(31*(i+1), len(data))]))
    }
    return poly
}

func BlobStoreRequestsToPoly(blobStoreRequests []BlobStoreRequest) ([]fr.Element, [][32]byte, []uint32) {
    blobIDs := make([][32]byte, 0)
    blobDataStartDegrees := make([]uint32, 0)
    overallPoly := make([]fr.Element, 0)
    for i := 0; i < len(blobStoreRequests); i++ {
        poly := DataToPoly(blobStoreRequests[i].Data)
        blobIDs = append(blobIDs, blobStoreRequests[i].BlobIDs)
        blobDataStartDegrees = append(blobDataStartDegrees, len(overallPoly))
        overallPoly = append(overallPoly, poly...)
    }
    return overallPoly, blobIDs, blobDataStartDegrees
}
```
The disperser returns to each requester the KZG commitment to the `overallPoly` that their data was included in, its start and end degrees, and the corresponding [DataStoreHeader](../spec/types/node-types.md#datastoreheader) that the blob was included in.

### Encoding

The disperser encodes the `overallPoly` for each quorum among all of the `BlobStoreRequests`. The disperser generates its encoding parameters for each quorum relative to the highest `AdversaryThresholdBPs` and highest `QuorumThresholdBPs` for each quorum among all of the `BlobStoreRequests`.

[TODO: @bxue-l2]

### Aggregation

## Confirmation
