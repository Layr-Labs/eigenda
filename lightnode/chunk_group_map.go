package lightnode

type ChunkGroupMap struct {

	// A map from light node ID to light node data.
	lightNodes map[uint64]LightNodeRegistration
}

// NewChunkGroupMap creates a new ChunkGroupMap.
func NewChunkGroupMap() ChunkGroupMap {
	return ChunkGroupMap{
		lightNodes: make(map[uint64]LightNodeRegistration),
	}
}

func (cgm *ChunkGroupMap) AddLightNode(lightNode LightNodeRegistration) {
	cgm.lightNodes[lightNode.ID] = lightNode
}
