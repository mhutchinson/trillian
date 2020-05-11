package main

import (
	"context"
	"crypto"
	"flag"
	"fmt"
	"time"

	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/x/beamx"
	"github.com/golang/glog"
	"github.com/google/trillian/experimental/beam/tilepb"
	"github.com/google/trillian/experimental/beam/tmap"
)

const hash = crypto.SHA512_256

var (
	valueSalt    = flag.String("value_salt", "v1", "Some string that will be smooshed in with the generated value before hashing. Allows updates for an index to be deterministic but variable.")
	keyCount     = flag.Int64("key_count", 1<<12, "Total number of key/value pairs to insert in the map.")
	inputShards  = flag.Int("input_shards", 2, "The number of shards to use to generate the input.")
	treeID       = flag.Int64("tree_id", 12345, "The ID of the tree. Used as a salt in hashing.")
	prefixStrata = flag.Int("prefix_strata", 3, "The number of strata of 8-bit strata before the final strata. 3 is optimal for trees up to 2^30. 10 is required to import into Trillian at the time of writing.")

	entries         = beam.NewCounter("demo", "entries")
	deNovoTiles     = beam.NewCounter("demo", "output-tiles")
	deNovoRootTiles = beam.NewCounter("demo", "output-root-tiles")
)

func main() {
	flag.Parse()
	beam.Init()
	defer glog.Flush()

	ctx := context.Background()
	p, s := beam.NewPipelineWithRoot()

	// Get a big old bag of leaves. This could read from an external source.
	leaves := getNewLeaves(s, *treeID, *keyCount, *inputShards, *valueSalt)

	// TODO(mhutchinson): Demo writing the output tiles to somewhere. Cloud Storage?
	allTiles, err := tmap.DeNovo(s, leaves, *treeID, hash, *prefixStrata)
	if err != nil {
		glog.Exitf("Failed to create pipeline: %v", err)
	}
	beam.ParDo(s, &CountFn{}, allTiles)

	before := time.Now()

	if err := beamx.Run(ctx, p); err != nil {
		glog.Exitf("Failed to execute job: %v", err)
	}
	glog.Infof("done in %v", time.Since(before).Round(time.Second))
}

func getNewLeaves(s beam.Scope, treeID, keyCount int64, shards int, valueSalt string) beam.PCollection {
	return beam.ParDo(s, &MapEntryFn{valueSalt, treeID}, CreateSequenceSource(s, keyCount, shards))
}

type MapEntryFn struct {
	Salt   string
	TreeID int64
}

func (fn *MapEntryFn) ProcessElement(ctx context.Context, i int64) *tilepb.MapEntry {
	h := hash.New()
	h.Write([]byte(fmt.Sprintf("%d", i)))
	kbs := h.Sum(nil)

	data := []byte(fmt.Sprintf("[%s]%d", fn.Salt, i))

	entries.Inc(ctx, 1)
	return tmap.CreateEntry(kbs, data, fn.TreeID)
}

// CreateSequenceSource creates a source PCollection of int64, which will be produced by
// a number of shards. This allows Flume to generate the sequence on multiple workers.
// The sequence is from [0, max).
func CreateSequenceSource(s beam.Scope, max int64, shards int) beam.PCollection {
	s = s.Scope("CreateSequenceSource")
	fullShardSize := max / int64(shards)

	var rangeStart int64
	ranges := make([]Range, shards)
	for i := 0; i < shards-1; i++ {
		end := rangeStart + fullShardSize
		ranges[i] = Range{rangeStart, end}
		rangeStart = end
	}
	ranges[shards-1] = Range{rangeStart, max}

	// This forces a shuffle operatiom to make the sharding happen. There's probably a better way.
	rangeToJunk := beam.GroupByKey(s, beam.SwapKV(s, beam.AddFixedKey(s, beam.CreateList(s, ranges))))
	return beam.ParDo(s, &EmitEntryFn{}, rangeToJunk)
}

type Range struct {
	Start, End int64
}

type EmitEntryFn struct {
}

func (fn *EmitEntryFn) ProcessElement(r Range, _ func(*int) bool, emit func(int64)) {
	for i := r.Start; i < r.End; i++ {
		emit(i)
	}
}

type CountFn struct {
}

func (fn *CountFn) ProcessElement(ctx context.Context, t *tilepb.MapTile) *tilepb.MapTile {
	deNovoTiles.Inc(ctx, 1)
	if len(t.GetIndex()) == 0 {
		deNovoRootTiles.Inc(ctx, 1)
	}
	return t
}
