# Integration Specific Glossary

EigenDA core glossary is located [here](https://docs.eigencloud.xyz/products/eigenda/core-concepts/glossary).

Rollup payload: compressed batches of transactions or state transition diffs.

DA cert: an `EigenDA Certificate` (or short `DACert`) contains all the information needed to retrieve a blob from the EigenDA network, as well as validate it.

EigenDA blob derivation: a sequence of procedures to convert a byte array repreesenting a DA cert to the final rollup payload.

Preimage oracle: Some procedures during the derivation require fetching additional data beyond the Input mentioned above. A preimage oracle is an object which the blob derivation is using for fetching those data.