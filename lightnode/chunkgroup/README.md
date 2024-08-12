# Chunk Groups

A "chunk group" is a collection of light nodes who are interested in sampling all chunks with a particular chunk index.
This document describes the algorithm for mapping each light node onto its chunk group.

# Desired Properties

The following properties are desired for the chunk group algorithm:

- **Randomness**: A light node operator should not be able to choose their chunk group prior to registration.
  A light node should have equal chance of being in any chunk particular chunk group.
- **Churn**: The chunk group of a light node should change over time.
- **Stability**: A light node should not change chunk groups too frequently, and at any point in time
  only a few light nodes should be be changing chunk groups.
- **Determinism**: The chunk group of a light node should be deterministic based on the node's  
  seed and the current time. Any two parties should agree which chunk group a light node is in
  at a particular timestamp.

# Terms

```
Consider the timeline below, with time moving from left to right.

                              The "+" marks represent                    The 7th time this
   The genesis time,          the time when a particular                 light node is shuffled.
   i.e. protocol start        light node is shuffled.                              |
          |                            |                                           ↓
          ↓      1          2          ↓          4          5          6          7          8          9
          |------+---|------+---|------+---|------+---|------+---|------+---|------+---|------+---|------+---|
          \          /          \      /                                \         /\         /\         /
           \        /            \    /                                  \       /  \       /  \       /
            \      /              \  /                                    \     /    \     /    \     /
             \    /                \/                                      \   /      \   /      \   /
              \  /                The "shuffle offset".                     \ /        \ /        \ /
               \/                 Each light node has a                   epoch 6     epoch 7    epoch 8
  A "shuffle period". Each node   random offset assigned
  changes chunk groups once per   at registration time.
  shuffle period. Each shuffle
  period is marked with a "|".
```

- **Genesis Time**: The time at which the protocol started. A light node's chunk group is only defined after the genesis
  time.
- **Shuffle Period**: The time that should pass in between a particular light node changing chunk groups (e.g. 1 week).
  Each time one shuffle period passes, all light nodes will be randomly reassigned to a new chunk group.
- **Shuffle Offset**: In order to avoid too many nodes changing chunk groups at the same time, each node switches chunk
  groups at an offset relative to the beginning of each shuffle period. This offset is randomly assigned to each node at
  registration time, but remains constant for the lifetime of the node.
- **Epoch**: An epoch describes the number of times a particular light node has changed chunk groups. At genesis,
  all light nodes are in epoch 0. The epoch for a particular light node is incremented by 1 for each time it is
  randomly shuffled into a new chunk group. The length of each epoch is equal to the shuffle period.

# Algorithm

The algorithm for determining which chunk group a particular light node is described below. A reference implementation
of this algorithm can be found in [calculations.go](./calculations.go).

## Determining a node's seed

A light node's seed is an 8 byte value that is randomly assigned to the node at registration time. The seed should
be generated using on-chain randomness that is difficult for an attacker to predict. The light node's seed
is public information and stored on-chain.

## `randomInt(seed)`

Define a function `randomInt(seed)` that takes an 8 byte signed integer as a seed and returns an 8 byte signed integer. 
TODO give spec for cryptographically secure random function!!!

## Determining a node's shuffle offset

A node's shuffle offset is a duration between 0 and the shuffle period at nanosecond granularity.
The shuffle offset is determined as follows:

```
nodeOffset_nanoseconds := randomInt(nodeSeed) % shufflePeriod_nanoseconds
```

## Epoch Calculation

At genesis, the epoch for each node is defined as `0`. The epoch increases to `1` at time equal
to `genesis + nodeOffset`, and then increases by one for each time the clock increases by a further `shufflePeriod`.

## Chunk Group Calculation

To determine a node's chunk group, first compute the current epoch for the node, then plug it into the following
function:

```
chunkGroup := randomInt(nodeSeed + nodeEpoch) % numberOfChunks
```