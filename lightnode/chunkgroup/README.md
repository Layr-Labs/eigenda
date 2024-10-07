# Chunk Groups

A "chunk group" is a collection of light nodes who are interested in sampling all chunks with a particular chunk index.
This document describes the algorithm for mapping each light node onto its chunk group(s).

# Desired Properties

The following properties are desired for the chunk assignment group algorithm:

- **Randomness**: A light node operator should not be able to choose their chunk group(s) prior to registration.
  A light node should have equal chance of being in any particular chunk group at any particular point in time.
- **Churn**: The chunk group of a light node should change over time.
- **Stability**: A light node should not change chunk groups too frequently, and at any point in time
  only a few light nodes should be be changing chunk groups.
- **Determinism**: The chunk group(s) of a light node should be deterministic based on the node's  
  seed and the current time. Any two parties should agree which chunk group(s) a light node is in
  at a particular timestamp.
- **Multiplicity**: This algorithm should support a light node being in 1 or more chunk groups simultaneously.
  (Note: a maximum of 64 chunk groups per light node is supported by this algorithm. The likely configuration for the
  number of chunk groups is 2.)

# Terms

```
Consider the timeline below, with time moving from left to right.

                              The "+" marks represent                    
                              the time when a particular                     The 7th shuffle
      unix epoch              light node shuffles a group.                         |
          |                            |                                           ↓
          ↓      1          2          ↓          4          5          6          7          8          9
          |------+---|------+---|------+---|------+---|------+---|------+---|------+---|------+---|------+---|
          \          /          \      /                                \         /\         /\         /
           \        /            \    /                                  \       /  \       /  \       /
            \      /              \  /                                    \     /    \     /    \     /
             \    /                \/                                      \   /      \   /      \   /
              \  /                The "shuffle offset".                     \ /        \ /        \ /
               \/                 For each group a light                  epoch 6     epoch 7    epoch 8
  A "shuffle period". Each node   nodes is in, it has a
  changes chunk groups once per   random offset assigned
  shuffle period per group it     at registration time.
  is in. Each shuffle period is 
  marked with a "|".
```

- **Unix Epoch**: time = 0, aka January 1, 1970. All times are measured in nanoseconds since the Unix epoch,
  and all light nodes considered to have an epoch number of 0 at precisely this time. Even though this protocol
  obviously did not exist at this time, it's convenient to use it as a reference point in order to keep the math simple.
- **Shuffle Period**: The time that should pass in between a particular light node changing one of its chunk groups 
  (e.g. 1 week). Each time one shuffle period passes, all light nodes within a chunk group will have been assigned to a 
  random chunk group (note that it's possible for a node to be randomly assigned to the same chunk group multiple 
  times in a row).
- **Shuffle Offset**: In order to avoid too many nodes changing chunk groups at the same time, each node switches chunk
  groups at an offset relative to the beginning of each shuffle period. This offset is randomly assigned to each node at
  registration time, but remains constant for the lifetime of the node. If a light node is in multiple chunk groups,
  then it will have a different shuffle offset for each group it is in.
- **Assignment Index**: a light node may participate in multiple chunk groups simultaneously. The first group it is in 
  is said to have an "assignment index" of 0, the second has an assignment index of 1, and so on. 
  The number of chunk groups that a light node participates in is a configuration parameter of the protocol.
  A light node shuffles each of its chunk groups independently. That is, the light node shuffles its first chunk group
  assignment at a random time that is independent of the time that it shuffles its second chunk group assignment 
  (and so on).
- **Shuffle Epoch**: An shuffle epoch describes the number of times a particular light node has changed chunk groups
  for a group assignment. At genesis, all light nodes' group assignments are in epoch 0. The epoch for a particular 
  light node's group assignment is incremented by 1 for each time it is randomly shuffled into a new chunk group. 
  The length of each epoch is equal to the shuffle period.


# Algorithm

The algorithm for determining which chunk group a particular light node is described below. A reference implementation
of this algorithm can be found in [calculations.go](./calculations.go).

## Determining a node's seed

A light node's seed is an 8 byte value that is randomly assigned to the node at registration time. The seed should
be generated using on-chain randomness that is difficult for an attacker to predict. The light node's seed
is public information and stored on-chain.

## `randomInt(uint64, uint64, ..., uint64)`

Define a function `randomInt(uint64, uint64, ..., uint64)` that takes a variable number of 8 byte unsigned integers
and returns a pseudo-random 8 byte unsigned integer. 

- For each unsigned integer from left to right, append the integer's bytes into a byte array called 
  `seedBytes` in big endian order.
- Use the `seedBytes` as the input to `keccak256` to generate a 32 byte array called `hashBytes`.
- Use the first 8 bytes of `hashBytes` to create an 8 byte unsigned integer using big endian order called `result`.
- Return `result`.

## Determining a node's shuffle offset for a particular assignment index

A node's shuffle offset for a particular assignment index is a duration between 0 and the shuffle period at 
nanosecond granularity. The shuffle offset is determined as follows:

```
shuffleOffset := randomInt(nodeSeed, assignmentIndex) % shufflePeriod_nanoseconds
```

## Shuffle Epoch Calculation

At the unix epoch, the shuffle epoch for each group assignment is defined as `0`. The epoch increases by `1` at time 
`shuffleOffset(nodeSeed, assignmentIndex)`, and then by `1` again every `shufflePeriod` nanoseconds.

## Chunk Group Calculation

To determine a node's chunk group, first compute the current epoch for the node, then plug it into the following
function:

```
chunkGroup := randomInt(nodeSeed, assignmentIndex, shuffleEpoch) % numberOfChunks
```