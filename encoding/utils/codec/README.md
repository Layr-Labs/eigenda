# How To Choose Good Payload Sizes

Choosing a good payload size is important for optimizing EigenDA usage costs. If you have the ability to control
the size of your payload and you choose un-optimally, you may end up paying twice as much for your traffic.

## Definitions

A `payload` is defined as the raw, unencoded data that is sent to EigenDA. From a logical point of view, this is what
an EigenDA customer wants to store and later be able to have that data be highly available.

A `blob` is a `payload` that has been encoded and packaged in a way that is suitable for sending to EigenDA. When a
`payload` is converted to a blob, the `blob` is always larger than the original `payload`. `Blobs` must always have
a length equal to a power of 2. If a `blob` would otherwise not be a power of 2, it is padded with zeros until its
length is a power of 2.

When EigenDA determines the cost of dispersing data, it uses the size of the `blob` as the basis for the cost, NOT
the size of the `payload`. If two `payloads` of different sizes are converted to a `blob` of the same size, they will
have the same cost. Since a `blob` size might be rounded up to the next power of 2, sometimes adding a single byte
to a `payload` can double the size of the resulting `blob`, and therefore double the cost of dispersing that data.

## Choosing Payload Sizes

The table below shows the `blob` size that various `payload` sizes will be converted to. Having a payload that exactly
matches a size in the `Maximum Payload Size` column means that the dispersal is maximally efficient from a cost
perspective. Going one byte over that size will double the cost of dispersing that data, as it pushes the `blob`
size to the next power of 2. If possible, aim to size your `payloads` to be as close to the `Maximum Payload Size` as
possible but without exceeding it.

In the table below, all bounds are inclusive.

| Maximum Payload Size | Blob Size               |
|:---------------------|:------------------------|
| 126945 bytes         | 131072 bytes (128 KiB)  |
| 253921 bytes         | 262144 bytes (256 KiB)  |
| 507873 bytes         | 524288 bytes (512 KiB)  |
| 1015777 bytes        | 1048576 bytes (1 MiB)   |
| 2031585 bytes        | 2097152 bytes (2 MiB)   |
| 4063201 bytes        | 4194304 bytes (4 MiB)   |
| 8126433 bytes        | 8388608 bytes (8 MiB)   |
| 16252897 bytes       | 16777216 bytes (16 MiB) |

## Minimum Blob SIze

The minimum `blob` size is 128KiB. Sending extremely small `payloads` will always result in being charted for at least
as much as if sending 128KiB. (Note that the actual data transmitted over the wire may be smaller than 128KiB, but
it is metered and charged as if it were 128KiB.)

## Maximum Blob Size

The maximum `blob` size is 16MiB. Sending extremely large `payloads` will result in a dispersal error if it cannot 
fit into a 16MiB `blob`.