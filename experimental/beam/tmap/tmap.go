package tmap

import (
	"crypto"
	"fmt"
	"reflect"
	"sort"

	"github.com/apache/beam/sdks/go/pkg/beam"

	"github.com/google/trillian/experimental/beam/tilepb"
	"github.com/google/trillian/merkle/coniks"
	"github.com/google/trillian/merkle/hashers"
	"github.com/google/trillian/merkle/smt"
	"github.com/google/trillian/storage/tree"
)

func init() {
	beam.RegisterCoder(reflect.TypeOf(tree.NodeID2{}), nodeID2Encode, nodeID2Decode)
}

// DeNovo creates a new map which is the result of converting
// the input PCollection of *tilepb.MapEntry into a Merkle tree and tiling the result.
// The output type is a PCollection of *tilepb.MapTile.
// prefixStrata is the number of 8-bit prefix strata. Any path from root to leaf will have $prefixStrata+1 tiles.
func DeNovo(s beam.Scope, leaves beam.PCollection, treeID int64, hash crypto.Hash, prefixStrata int) (beam.PCollection, error) {
	s = s.Scope("tmap.DeNovo")
	if prefixStrata < 0 || prefixStrata >= 32 {
		return beam.Impulse(s), fmt.Errorf("prefixStrata must be in [0, 32), got %d", prefixStrata)
	}

	// Construct the map pipeline starting with the leaf tile strata.
	// coordinateHeight is the height of the root of the current strata being constructed.
	coordinateHeight := 256 - (8 * prefixStrata)
	lShard, lAgg := createLeafFns(treeID, hash, coordinateHeight)
	lastStrata := beam.ParDo(s, lAgg, beam.GroupByKey(s, beam.ParDo(s, lShard, leaves)))
	allTiles := []beam.PCollection{lastStrata}
	for i := 0; i < prefixStrata; i++ {
		coordinateHeight += 8
		// TODO(mhutchinson): Convert all units to bytes in this class; bits should go away.
		tShard, tAgg := createTileFns(treeID, hash, 8, coordinateHeight)
		lastStrata = beam.ParDo(s, tAgg, beam.GroupByKey(s, beam.ParDo(s, tShard, lastStrata)))
		allTiles = append(allTiles, lastStrata)
	}

	// Collate all of the strata together and return them.
	return beam.Flatten(s, allTiles...), nil
}

// Update takes an existing base map (PCollection<*tilepb.MapTile>), applies the
// delta (PCollection<*tilepb.MapEntry>) and returns the tiled map (PCollection<*tilepb.MapTile>).
// The treeID & strata of the base map must match the scheme passed in, or things will break. Horribly.
func Update(s beam.Scope, base, delta beam.PCollection, treeID int64, hash crypto.Hash, prefixStrata int) (beam.PCollection, error) {
	s = s.Scope("tmap.Update")
	if prefixStrata < 0 || prefixStrata >= 32 {
		return beam.Impulse(s), fmt.Errorf("prefixStrata must be in [0, 32), got %d", prefixStrata)
	}

	// Construct the map pipeline starting with the leaf tile strata.
	// coordinateHeight is the height of the root of the current strata being constructed.
	coordinateHeight := 256 - (8 * prefixStrata)
	shardedLeafBase := beam.ParDo(s, &BaseTileSharder{(256 - coordinateHeight) / 8}, base)
	shardedLeafDelta := beam.ParDo(s, &LeafShardFn{uint(256 - coordinateHeight)}, delta)
	lastStrata := beam.ParDo(s, &LeafTileUpdateFn{TreeID: treeID, Hash: hash, HeightBytes: coordinateHeight / 8}, beam.CoGroupByKey(s, shardedLeafBase, shardedLeafDelta))
	allTiles := []beam.PCollection{lastStrata}
	for i := 0; i < prefixStrata; i++ {
		coordinateHeight += 8
		shardedTileBase := beam.ParDo(s, &BaseTileSharder{(256 - coordinateHeight) / 8}, base)
		shardedTileDelta := beam.ParDo(s, &TileShardFn{uint(256 - coordinateHeight)}, lastStrata)
		lastStrata = beam.ParDo(s, &UpperTileUpdateFn{TreeID: treeID, Hash: hash, HeightBytes: 1, RootHeight: coordinateHeight}, beam.CoGroupByKey(s, shardedTileBase, shardedTileDelta))
		allTiles = append(allTiles, lastStrata)
	}

	// Collate all of the strata together and return them.
	return beam.Flatten(s, allTiles...), nil
}

// CreateEntry creates an entry given the 256 bit key and the *raw value* that this should map
// to. The raw value can be any length and will be hashed. This would ideally be changed so
// that Trillian doesn't *need* to ever see raw values. However, the current Trillian Map API and
// verification library require that the raw value is available. The ideal design for the map is
// one in which the map only commits to a dictionary of hash->hash, and the personality can provide
// layers on top of that that provide the value and the strategy for turning that into the value hash.
func CreateEntry(key, rawData []byte, treeID int64) *tilepb.MapEntry {
	vbs := coniks.Default.HashLeaf(treeID, key, rawData)

	return &tilepb.MapEntry{
		Key:   key,
		Value: vbs,
		Data:  rawData,
	}
}

func createLeafFns(treeID int64, hash crypto.Hash, tileHeightBits int) (*LeafShardFn, *LeafAggregatorFn) {
	return &LeafShardFn{uint(256 - tileHeightBits)}, &LeafAggregatorFn{TreeID: treeID, Hash: hash, HeightBytes: tileHeightBits / 8}
}

func createTileFns(treeID int64, hash crypto.Hash, tileHeightBits, coordinateHeightBits int) (*TileShardFn, *TileAggregatorFn) {
	return &TileShardFn{uint(256 - coordinateHeightBits)}, &TileAggregatorFn{TreeID: treeID, Hash: hash, HeightBytes: tileHeightBits / 8, RootHeight: coordinateHeightBits}
}

type TileShardFn struct {
	// The depth of the root of this tile from the root of the tree.
	RootDepth uint
}

func (fn *TileShardFn) ProcessElement(t *tilepb.MapTile) (tree.NodeID2, *tilepb.MapTile) {
	p := t.GetIndex()
	return tree.NewNodeID2(string(p), fn.RootDepth), t
}

type TileAggregatorFn struct {
	HeightBytes int
	// The height of the root of this tile above the lowest level of the tree. This is bits, not bytes.
	RootHeight int
	TreeID     int64
	Hash       crypto.Hash
	th         *tileHasher
}

func (fn *TileAggregatorFn) Setup() {
	fn.th = &tileHasher{fn.TreeID, coniks.New(fn.Hash)}
}

func (fn *TileAggregatorFn) ProcessElement(root tree.NodeID2, leaves func(**tilepb.MapTile) bool) (*tilepb.MapTile, error) {
	m := make(map[tree.NodeID2]leafValue)
	var leaf *tilepb.MapTile
	for leaves(&leaf) {
		lidx := tree.NewNodeID2(string(leaf.GetIndex()), root.BitLen()+8*uint(fn.HeightBytes))
		if v, found := m[lidx]; found {
			return nil, fmt.Errorf("Found duplicate values at tile position %s: {%x, %x}", lidx, v, leaf.GetRoot())
		}
		m[lidx] = leafValue{leaf.GetRoot(), nil}
	}

	return fn.th.createTile(8*fn.HeightBytes, root, m)
}

// LeafShardFn is a DoFn that shards the entries by the node id of the root of the tile
// that will include them.
type LeafShardFn struct {
	// The height of the root of this tile above the lowest level of the tree. This is bits, not bytes.
	RootDepth uint
}

func (fn *LeafShardFn) ProcessElement(leaf *tilepb.MapEntry) (tree.NodeID2, *tilepb.MapEntry) {
	return tree.NewNodeID2(string(leaf.GetKey()), fn.RootDepth), leaf
}

type LeafAggregatorFn struct {
	HeightBytes int
	TreeID      int64
	Hash        crypto.Hash
	th          *tileHasher
}

func (fn *LeafAggregatorFn) Setup() {
	fn.th = &tileHasher{fn.TreeID, coniks.New(fn.Hash)}
}

func (fn *LeafAggregatorFn) ProcessElement(root tree.NodeID2, leaves func(**tilepb.MapEntry) bool) (*tilepb.MapTile, error) {
	m := make(map[tree.NodeID2]leafValue)
	var leaf *tilepb.MapEntry
	for leaves(&leaf) {
		lidx := tree.NewNodeID2(string(leaf.GetKey()), root.BitLen()+8*uint(fn.HeightBytes))
		if v, found := m[lidx]; found {
			return nil, fmt.Errorf("Found duplicate values at leaf tile position %s: {%x, %x}", lidx, v, leaf.GetValue())
		}
		m[lidx] = leafValue{leaf.GetValue(), leaf.GetData()}
	}

	rootHeight := 8 * fn.HeightBytes
	return fn.th.createTile(rootHeight, root, m)
}

// BaseTileSharder keeps only the tiles from the input PCollection that are at the given
// depth. It will then shard these tiles by the NodeID2 of the root.
// This could be rewritten as a filter.Include and then TileShardFn.
type BaseTileSharder struct {
	DepthBytes int
}

func (fn *BaseTileSharder) ProcessElement(t *tilepb.MapTile, emit func(tree.NodeID2, *tilepb.MapTile)) {
	if len(t.Index) == fn.DepthBytes {
		node, _ := nodeID2Decode(t.Index)
		emit(node, t)
	}
}

// UpperTileUpdateFn takes care of updates for tiles in the upper strata.
// The bases input is 0 or 1 tiles which is the old state of this tile.
// The leaves input is any number of leaf tiles that are the new children of this.
type UpperTileUpdateFn struct {
	HeightBytes int
	// The height of the root of this tile above the lowest level of the tree. This is bits, not bytes.
	RootHeight int
	TreeID     int64
	Hash       crypto.Hash
	th         *tileHasher
}

func (fn *UpperTileUpdateFn) Setup() {
	fn.th = &tileHasher{fn.TreeID, coniks.New(fn.Hash)}
}

func (fn *UpperTileUpdateFn) ProcessElement(root tree.NodeID2, bases func(**tilepb.MapTile) bool, leaves func(**tilepb.MapTile) bool) (*tilepb.MapTile, error) {
	base, err := getZeroOrOne(bases)
	if err != nil {
		return nil, fmt.Errorf("Error at location %s: %v", root, err)
	}
	// Consider short-circuiting in the case that leaves is empty.

	m := make(map[tree.NodeID2]leafValue)
	// 1. Init the map with the old values for each leaf in this tile.
	if base != nil {
		rootBytes, err := nodeID2Encode(root)
		if err != nil {
			return nil, err
		}
		for _, l := range base.Leaves {
			lidx, err := nodeID2Decode(append(rootBytes, l.Path...))
			if err != nil {
				return nil, err
			}
			m[lidx] = leafValue{l.GetValue(), nil}
		}
	}

	// 2. Write and new values into the map, overwriting any stale entries.
	var leaf *tilepb.MapTile
	for leaves(&leaf) {
		lidx := tree.NewNodeID2(string(leaf.GetIndex()), root.BitLen()+8*uint(fn.HeightBytes))
		m[lidx] = leafValue{leaf.GetRoot(), nil}
	}
	return fn.th.createTile(8*fn.HeightBytes, root, m)
}

type LeafTileUpdateFn struct {
	HeightBytes int
	TreeID      int64
	Hash        crypto.Hash
	th          *tileHasher
}

func (fn *LeafTileUpdateFn) Setup() {
	fn.th = &tileHasher{fn.TreeID, coniks.New(fn.Hash)}
}

func (fn *LeafTileUpdateFn) ProcessElement(root tree.NodeID2, bases func(**tilepb.MapTile) bool, deltas func(**tilepb.MapEntry) bool) (*tilepb.MapTile, error) {
	base, err := getZeroOrOne(bases)
	if err != nil {
		return nil, fmt.Errorf("Error at location %s: %v", root, err)
	}
	noDeltas, deltas := isEmpty(deltas)

	if noDeltas && base == nil {
		return nil, fmt.Errorf("noDeltas and no base should imply no function call.")
	}
	if noDeltas {
		// If there are no deltas, then the base tile is unchanged.
		return base, nil
	}

	// We add new values first and then update with base to easily check for duplicates in deltas.
	m := make(map[tree.NodeID2]leafValue)
	var leaf *tilepb.MapEntry
	for deltas(&leaf) {
		lidx := tree.NewNodeID2(string(leaf.GetKey()), root.BitLen()+8*uint(fn.HeightBytes))
		if v, found := m[lidx]; found {
			return nil, fmt.Errorf("Found duplicate values at leaf tile position %s: {%x, %x}", lidx, v, leaf.GetValue())
		}
		m[lidx] = leafValue{leaf.GetValue(), leaf.GetData()}
	}

	if base != nil {
		rootBytes, err := nodeID2Encode(root)
		if err != nil {
			return nil, err
		}
		for _, l := range base.Leaves {
			lidx, err := nodeID2Decode(append(rootBytes, l.Path...))
			if err != nil {
				return nil, err
			}
			if _, found := m[lidx]; !found {
				// Only add base values if they haven't been updated.
				m[lidx] = leafValue{l.GetValue(), l.GetData()}
			}
		}
	}

	rootHeight := 8 * fn.HeightBytes
	return fn.th.createTile(rootHeight, root, m)
}

type leafValue struct {
	value, data []byte
}

// tileHasher is a NodeAccessor used for computing node hashes of a tile.
type tileHasher struct {
	treeID int64
	h      hashers.MapHasher
}

func (th *tileHasher) createTile(tileHeight int, root tree.NodeID2, lm map[tree.NodeID2]leafValue) (*tilepb.MapTile, error) {
	leafNodes := make([]smt.Node, len(lm))
	leaves := make([]*tilepb.TileLeaf, len(lm))
	i := 0
	for k := range lm {
		lv := lm[k]
		leafNodes[i] = smt.Node{k, lv.value}

		bs, _ := nodeID2Encode(k)
		leaves[i] = &tilepb.TileLeaf{
			Path:  bs[root.BitLen()/8:],
			Value: lv.value,
			Data:  lv.data,
		}
		i++
	}
	rootHash, err := th.hashTile(root, leafNodes)
	if err != nil {
		return nil, err
	}
	sort.Sort(byLeafPath(leaves))

	indexBytes, err := nodeID2Encode(root)
	if err != nil {
		return nil, err
	}
	return &tilepb.MapTile{
		Index:      indexBytes,
		TileLevels: int32(tileHeight),
		Root:       rootHash,
		Leaves:     leaves,
	}, nil
}

func (th *tileHasher) hashTile(root tree.NodeID2, leaves []smt.Node) ([]byte, error) {
	h, err := smt.NewHStar3(leaves, th.h.HashChildren, uint(leaves[0].ID.BitLen()), root.BitLen())
	if err != nil {
		return nil, err
	}
	r, err := h.Update(th)
	if err != nil {
		return nil, err
	}
	if len(r) != 1 {
		return nil, fmt.Errorf("Expected single root but got %d", len(r))
	}
	return r[0].Hash, nil
}

// Get returns an empty hash for the given root node ID.
func (th tileHasher) Get(id tree.NodeID2) ([]byte, error) {
	oldID := tree.NewNodeIDFromID2(id)
	height := th.h.BitLen() - oldID.PrefixLenBits
	return th.h.HashEmpty(th.treeID, oldID.Path, height), nil
}

func (th tileHasher) Set(id tree.NodeID2, hash []byte) {}

type byLeafPath []*tilepb.TileLeaf

func (a byLeafPath) Len() int      { return len(a) }
func (a byLeafPath) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byLeafPath) Less(i, j int) bool {
	iPath, jPath := a[i].Path, a[j].Path
	for bi := 0; bi < len(iPath); bi++ {
		if iPath[bi] != jPath[bi] {
			return iPath[bi] < jPath[bi]
		}
	}
	return false
}

func getZeroOrOne(es func(**tilepb.MapTile) bool) (*tilepb.MapTile, error) {
	var r, e *tilepb.MapTile
	for es(&e) {
		if r != nil {
			return nil, fmt.Errorf("Expected zero or one entry, found at least two (%s, %s)", r, e)
		}
		r = e
	}
	return r, nil
}

// isEmpty returns true if the iterator has no elements, and also returns the iterator
// function that behaves as the original iterator function did. Clients should discard
// the original iterator after passing into this function.
func isEmpty(es func(**tilepb.MapEntry) bool) (bool, func(**tilepb.MapEntry) bool) {
	var next *tilepb.MapEntry
	hasNext := es(&next)
	if !hasNext {
		return true, es
	}
	// Now we're stuck in this world where we're always an element behind.
	return false, func(r **tilepb.MapEntry) bool {
		if !hasNext {
			return false
		}
		*r = next
		hasNext = es(&next)
		return true
	}
}

func nodeID2Encode(n tree.NodeID2) ([]byte, error) {
	b, c := n.LastByte()
	if c%8 != 0 {
		return nil, fmt.Errorf("Found non-byte-aligned node id with bit length %d", n.BitLen())
	}
	if c == 0 {
		return []byte{}, nil
	}
	return append([]byte(n.FullBytes()), b), nil
}

func nodeID2Decode(bs []byte) (tree.NodeID2, error) {
	return tree.NewNodeID2(string(bs), 8*uint(len(bs))), nil
}
