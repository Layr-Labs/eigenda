
# Indexer

## Base Indexer

### Description

The indexer has the function of maintaining accumulators based on events.

An accumulator consists of an object, an event filter, and a mutator. Whenever an event matching the filter is received, the event is fed to the mutator, which updates the object (possibly making calls to the current smart contract state). We can configure how much history of the accumulator to store. 

An accumulator can optionally define a custom method for initializing the accumulator from state from an intermediate checkpoint. This includes methods such as pulling state from a smart contract or getting calldata associated with a particular transaction. 

The indexer is one of the only stateful components of the operator. To avoid reindexing on restarts, the state of the indexer is stored in a database. We will use a schemaless db to avoid migrations.

The indexer must also support reorg resistance. We can achieve simple reorg resilience in the following way:
- For every accumulator, we make sure to store history long enough that we always have access to a finalized state. 
- In the event reorg is detected, we can revert to the most recent finalized state, and then reindex to head. 

The indexer needs to accommodate upgrades to the smart contract interfaces. Contract upgrades can have the following effects on interfaces:
- Addition, removal, modification of events
- Addition, removal, modification of state variables

The indexer has native support for "forking" to support upgrades to contract interfaces. A fork consists of the following:
- A forking condition
- A new set of accumulator implementations; each accumulator implements a fork method which may implement some update 

For each accumulator, the indexer exposes an interface, which can be queried by block number (or perhaps other fields). The returned response includes the accumulator object as well as the fork value. 

### Spec


```go

type Indexer struct{
    Accumulators []*Accumulator
    HeaderService *HeaderService
    HeaderStore *HeaderStore
    UpgradeForkWatcher *UpgradeForkWatcher
}


const maxUint := ^uint(0)
const maxSyncBlocks = 10

func (i Indexer) Index(){

    // Find the latest block that we can fast forward to. 

    clientLatestHeader := i.HeaderService.PullLatestHeader(true)

    syncFromBlock := maxUint

    for _, acc := range i.Accumulators{
        bn := acc.GetSyncPoint(clientLatestHeader)
        if syncFromBlock > bn{
            syncFromBlock = bn
        }
    }

    bn := i.UpgradeForkWatcher.GetLatestUpgrade(clientLatestHeader)
    if syncFromBlock > bn{
        syncFromBlock = bn
    }

    myLatestHeader := i.HeaderService.GetLatestHeader(true)

    // TODO: Also if there are no headers synced
    // Fast forward if it's too many blocks to catch up
    if syncFromBlock - myLatestHeader.Number > maxSyncBlocks {

        // This probably just wipes the HeaderStore clean
        i.HeaderStore.FastForward()

        for _, acc := range i.Accumulators{
            acc.SetSyncPoint(clientLatestHeader)
        }
    }


    go func(){
        for {
            headers, _ := i.HeaderService.PullNewHeaders()

            headers := i.UpgradeForkWatcher.DetectUpgrade(headers)

            newHeaders, _ := i.HeaderStore.AddHeaders(headers)

            for _, acc := range i.Accumulators{
                i.HandleAccumulator(acc,newHeaders)
            }
        }
    }
}

func (i Indexer) HandleAccumulator(acc Accumulator, headers []Headers){

    // Handle fast mode
    header, object, headers, err := acc.HandleFastMode(headers)
    if object != nil{
        i.HeaderStore.AttachObject(header, object)
    }

    // Get the starting accumulator object
    object := i.HeaderStore.GetLatestObject(acc)

    // Process headers
    object, headers := acc.ProcessHeaders(object, headers)

    // Register these accumulator objects
    for ind = range objects{
        i.HeaderStore.AttachObject(headers[i], acc, objects[i])
    }

}

type UpgradeFork struct {
    name string
} 

type Header struct {
    BlockHash [32]byte
    PrevBlockHash [32]byte
    Number uint64
    Finalized bool
    CurrentFork *UpgradeFork
    IsUpgrade bool
}


// UpgradeForkWatcher is a component that is used to scan a list of headers for an upgrade. Future upgrades may be based on a condition; past upgrades should have a block number configuration provided. 
type UpgradeForkWatcher() interface{

    // DetectUpgrade takes in a list of headers and sets the CurrentFork and IsUpgrade fields
    DetectUpgrade(headers []Header) []Header

}

// HeaderStore is a stateful component that maintains a chain of headers and their finalization status.
type HeaderStore interface{

    // Addheaders finds the header  It then crawls along this list of headers until it finds the point of divergence with its existing chain. All new headers from this point of divergence onward are returned.
    AddHeaders(headers []Headers) ([]Header, error)

    // GetLatestHeader returns the most recent header that the HeaderService has previously pulled
    GetLatestHeader(finalized bool)

    // AttachObject takes an accumulator object and attaches it to a header so that it can be retrieved using GetObject
    AttachObject(header Header, acc Accumulator, object AccumulatorObject) error

    // GetObject takes in a header and retrieves the accumulator object attached to the latest header prior to the supplied header having the requested object type. 
    GetObject(header Header, acc Accumulator) (AccumulatorObject, Header, error)

    // GetObject retrieves the accumulator object attached to the latest header having the requested object type. 
    GetLatestObject(acc Accumulator) (AccumulatorObject, Header, error)

}

// HeaderService
type HeaderService interface{

    // GetHeaders returns a list of new headers since the indicated header.  PullNewHeaders automatically handles batching and waiting for a specified period if it is already at head. 
    PullNewHeaders(lastHeader Header) ([]Header, error)

    // PullLatestHeader gets the latest header from the chain client
    PullLatestHeader(finalized bool)

}

type AccumulatorObject interface{
    
}

type Accumulator interface{

    // IndexHeaders accepts a list of incoming headers. Will throw an error is the accumulator does not have an existing header which can form a chain with the incoming headers. The Accumulator will discard any orphaned headers. 
    ProcessHeaders(headers []Headers) error

    // GetAccumulator accepts a header and returns the value of the accumulator at that header. Either the Number or BlockHash fields of the header can be used. 
    GetAccumulator(header Header) (interface{},UpgradeFork,error)

    // GetAccumulators extends GetAccumulator to multiple blocks
    GetAccumulators(header Header) (interface{},UpgradeFork,error)

    // GetSyncPoint determines the blockNumber at which it needs to start syncing from based on both 1) its ability to full its entire state from the chain and 2) its indexing duration requirements.
    GetSyncPoint(latestHeader Header) (uint64,error)

    // SetSyncPoint sets the Accumulator to operate in fast mode. 
    SetSyncPoint(latestHeader Header) (error)

    // HandleFastMode handles the fast mode operation of the accumulator. In this mode, it will ignore all headers until it reaching the blockNumber associated with GetSyncPoint. Upon reaching this blockNumber, it will pull its entire state from the chain and then proceed with normal syncing. 
    HandleFastMode(latestHeader Header) (header, *AccumulatorObject, []header, error)

}

```

Reference accumulator implementation:

```go

type StateAccumulator struct{
    
    HistoryLength uint64
    IsFastMode bool
    EndFastModeBlockNumber uint64
}

type StateObject interface{

}

type StateObjectV1 struct{

}

func (a *StateAccumulator) GetSyncPoint(latestHeader Header) uint64{
    return latestHeader.Number - a.HistoryLength
}

func (a *StateAccumulator) SetSyncPoint(latestHeader Header) error{
    a.EndFastModeBlockNumber = a.GetSyncPoint(latestHeader)
    a.IsFastMode = true
}

func (a *StateAccumulator) HandleFastMode(latestHeader Header) ([]header, *AccumulatorObject, error){

    // Do fastMode processing
    if a.IsFastMode{

        for ind, header := range headers{
            
            if header.Number >= a.EndFastModeBlockNumber{
                object, err := a.InitializeStateObject(header)

                a.IsFastMode = false
                
                break
            }
        }
        headers = headers[ind+1:]
    }
    return headers, nil

}


func (a *StateAccumulator) InitializeStateObject(header Header) (StateObject,error) {
}

func (a *StateAccumulator) UpdateStateObject() (StateObject,error) {
}

func (a *StateAccumulator) ProcessHeaders(headers []Headers){
}


func (a *StateAccumulator) GetAccumulator(header Header) (interface{},UpgradeFork,error){
}

func GetAccumulators(header Header) (interface{},UpgradeFork,error){

}




```


## DA Indexer

The indexer keeps track of the [state variables](./contracts.md#bls-registry) from the `BLSRegistry` contract. 


