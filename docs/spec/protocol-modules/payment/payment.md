
## Payment Guarantee

THIS SECTION IS A DRAFT.

**Considerations for how payment verification affects indexing requirements**

- If I need to be able to calculate everyone's payment, then I need to be able to index back to an arbitrarily early state.
- If I only need to be able to calculate my own payment, then I only need to be able to index from my starting point. I can also do accumulation; don't need to store everything. 
- If cumulative payments are solidified on chain periodically, then I don't have to index back so far. 

Posting cumulative payments as calldata could be somewhat expensive; we would only want to do this one per period. 

Preliminary decision: We can build the indexer to 1) sync from the last time that cumulative payments were posted, 2) act as an accumulator so that it doesn't have to store old events. Then we can separately design ways to incentivise/enforce periodic posting of cumulative payments on chain. 

What actually needs to be stored in order to calculate payments? All of the staking amounts are in the persistent on-chain state. All that is needed is the sizes of the blobs, payment rates, etc.
