package tmap

import (
	"crypto"
	"fmt"
	"math/rand"
	"testing"

	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/testing/passert"
	"github.com/apache/beam/sdks/go/pkg/beam/testing/ptest"
	"github.com/apache/beam/sdks/go/pkg/beam/transforms/filter"

	"github.com/google/trillian/experimental/beam/tilepb"
)

const treeID = int64(12345)
const hash = crypto.SHA512_256

func TestMain(m *testing.M) {
	ptest.Main(m)
}

func TestDeNovo(t *testing.T) {
	tests := []struct {
		prefixStrata int
		entries      []*tilepb.MapEntry
		treeID       int64
		hash         crypto.Hash
		root         string

		expectFailConstruct bool
		expectFailRun       bool
	}{
		{
			prefixStrata: 0,
			entries:      []*tilepb.MapEntry{createEntry("ak", "av")},
			treeID:       12345,
			hash:         crypto.SHA512_256,
			root:         "af079c268bd48eb89532b2b0c96d753c8f98eb8ce03f5dd95fa60ab9cc92f3a4",
		},
		{
			prefixStrata: 0,
			entries:      []*tilepb.MapEntry{createEntry("ak", "av")},
			treeID:       54321,
			hash:         crypto.SHA512_256,
			root:         "8e6363380169b790b6e3d1890fc3d492a73512d9bbbfb886854e10ca10fc147f",
		},
		{
			prefixStrata: 1,
			entries:      []*tilepb.MapEntry{createEntry("ak", "av")},
			treeID:       12345,
			hash:         crypto.SHA512_256,
			root:         "af079c268bd48eb89532b2b0c96d753c8f98eb8ce03f5dd95fa60ab9cc92f3a4",
		},
		{
			prefixStrata: 0,
			entries:      []*tilepb.MapEntry{createEntry("ak", "av"), createEntry("bk", "bv"), createEntry("ck", "cv")},
			treeID:       12345,
			hash:         crypto.SHA512_256,
			root:         "2372f0432e04dc76015f427ce8a1294644e36421b047ddfd52afdfdba60aff25",
		},
		{
			prefixStrata: 1,
			entries:      []*tilepb.MapEntry{createEntry("ak", "av"), createEntry("bk", "bv"), createEntry("ck", "cv")},
			treeID:       12345,
			hash:         crypto.SHA512_256,
			root:         "2372f0432e04dc76015f427ce8a1294644e36421b047ddfd52afdfdba60aff25",
		},
		{
			prefixStrata:  0,
			entries:       []*tilepb.MapEntry{createEntry("ak", "av"), createEntry("ak", "av")},
			treeID:        12345,
			hash:          crypto.SHA512_256,
			expectFailRun: true,
		},
		{
			prefixStrata:        -1,
			entries:             []*tilepb.MapEntry{createEntry("ak", "av")},
			treeID:              12345,
			hash:                crypto.SHA512_256,
			expectFailConstruct: true,
		},
		{
			prefixStrata:        32,
			entries:             []*tilepb.MapEntry{createEntry("ak", "av")},
			treeID:              12345,
			hash:                crypto.SHA512_256,
			expectFailConstruct: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run("", func(t *testing.T) {
			// t.Parallel() would be great, but it seems Flume tests are sketchy about this?
			p, s := beam.NewPipelineWithRoot()
			leaves := beam.CreateList(s, test.entries)

			tiles, err := DeNovo(s, leaves, test.treeID, test.hash, test.prefixStrata)
			if got, want := err != nil, test.expectFailConstruct; got != want {
				t.Errorf("pipeline construction failure: got %v, want %v (%v)", got, want, err)
			}
			if test.expectFailConstruct {
				return
			}
			rootTile := filter.Include(s, tiles, func(t *tilepb.MapTile) bool { return len(t.GetIndex()) == 0 })
			roots := beam.ParDo(s, func(t *tilepb.MapTile) string { return fmt.Sprintf("%x", t.Root) }, rootTile)

			passert.Equals(s, roots, test.root)
			err = ptest.Run(p)
			if got, want := err != nil, test.expectFailRun; got != want {
				t.Errorf("pipeline run failure: got %v, want %v (%v)", got, want, err)
			}
		})
	}
}

func TestUpdateEmptyBase(t *testing.T) {
	tests := []struct {
		prefixStrata int
		entries      []*tilepb.MapEntry
		root         string

		expectFailConstruct bool
		expectFailRun       bool
	}{
		{
			prefixStrata: 0,
			entries:      []*tilepb.MapEntry{createEntry("ak", "av")},
			root:         "af079c268bd48eb89532b2b0c96d753c8f98eb8ce03f5dd95fa60ab9cc92f3a4",
		},
		{
			prefixStrata: 1,
			entries:      []*tilepb.MapEntry{createEntry("ak", "av")},
			root:         "af079c268bd48eb89532b2b0c96d753c8f98eb8ce03f5dd95fa60ab9cc92f3a4",
		},
		{
			prefixStrata: 0,
			entries:      []*tilepb.MapEntry{createEntry("ak", "av"), createEntry("bk", "bv"), createEntry("ck", "cv")},
			root:         "2372f0432e04dc76015f427ce8a1294644e36421b047ddfd52afdfdba60aff25",
		},
		{
			prefixStrata: 1,
			entries:      []*tilepb.MapEntry{createEntry("ak", "av"), createEntry("bk", "bv"), createEntry("ck", "cv")},
			root:         "2372f0432e04dc76015f427ce8a1294644e36421b047ddfd52afdfdba60aff25",
		},
	}

	for _, test := range tests {
		test := test
		t.Run("", func(t *testing.T) {
			// t.Parallel() would be great, but it seems Flume tests are sketchy about this?
			p, s := beam.NewPipelineWithRoot()

			// This actually puts an invalid MapTile into the base PCollection because creating an
			// empty PCollection seems impossible to cleanly do.
			base := beam.CreateList(s, []*tilepb.MapTile{{Index: []byte{0xff, 0xff}}})
			delta := beam.CreateList(s, test.entries)

			tiles, err := Update(s, base, delta, treeID, hash, test.prefixStrata)
			if got, want := err != nil, test.expectFailConstruct; got != want {
				t.Errorf("pipeline construction failure: got %v, want %v (%v)", got, want, err)
			}
			if test.expectFailConstruct {
				return
			}
			rootTile := filter.Include(s, tiles, func(t *tilepb.MapTile) bool { return len(t.GetIndex()) == 0 })
			roots := beam.ParDo(s, func(t *tilepb.MapTile) string { return fmt.Sprintf("%x", t.Root) }, rootTile)

			passert.Equals(s, roots, test.root)
			err = ptest.Run(p)
			if got, want := err != nil, test.expectFailRun; got != want {
				t.Errorf("pipeline run failure: got %v, want %v (%v)", got, want, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		prefixStrata               int
		baseEntries, updateEntries []*tilepb.MapEntry
		root                       string
	}{
		{
			prefixStrata:  0,
			baseEntries:   []*tilepb.MapEntry{createEntry("ak", "ignored")},
			updateEntries: []*tilepb.MapEntry{createEntry("ak", "av")},
			root:          "af079c268bd48eb89532b2b0c96d753c8f98eb8ce03f5dd95fa60ab9cc92f3a4",
		},
		{
			prefixStrata:  1,
			baseEntries:   []*tilepb.MapEntry{createEntry("ak", "ignored")},
			updateEntries: []*tilepb.MapEntry{createEntry("ak", "av")},
			root:          "af079c268bd48eb89532b2b0c96d753c8f98eb8ce03f5dd95fa60ab9cc92f3a4",
		},
		{
			prefixStrata:  0,
			baseEntries:   []*tilepb.MapEntry{createEntry("ak", "ignored"), createEntry("bk", "bv")},
			updateEntries: []*tilepb.MapEntry{createEntry("ak", "av"), createEntry("ck", "cv")},
			root:          "2372f0432e04dc76015f427ce8a1294644e36421b047ddfd52afdfdba60aff25",
		},
		{
			prefixStrata:  1,
			baseEntries:   []*tilepb.MapEntry{createEntry("ak", "ignored"), createEntry("bk", "bv")},
			updateEntries: []*tilepb.MapEntry{createEntry("ak", "av"), createEntry("ck", "cv")},
			root:          "2372f0432e04dc76015f427ce8a1294644e36421b047ddfd52afdfdba60aff25",
		},
	}

	for _, test := range tests {
		test := test
		t.Run("", func(t *testing.T) {
			// t.Parallel() would be great, but it seems Flume tests are sketchy about this?
			p, s := beam.NewPipelineWithRoot()
			leaves := beam.CreateList(s, test.baseEntries)

			base, err := DeNovo(s, leaves, treeID, hash, test.prefixStrata)
			if err != nil {
				t.Fatalf("failed to create pipeline: %v", err)
			}

			delta := beam.CreateList(s, test.updateEntries)
			update, err := Update(s, base, delta, treeID, hash, test.prefixStrata)
			if err != nil {
				t.Fatalf("failed to create pipeline: %v", err)
			}
			rootTile := filter.Include(s, update, func(t *tilepb.MapTile) bool { return len(t.GetIndex()) == 0 })
			roots := beam.ParDo(s, func(t *tilepb.MapTile) string { return fmt.Sprintf("%x", t.Root) }, rootTile)

			passert.Equals(s, roots, test.root)
			if err := ptest.Run(p); err != nil {
				t.Fatalf("pipeline failed: %v", err)
			}
		})
	}
}

func TestChildrenSorted(t *testing.T) {
	p, s := beam.NewPipelineWithRoot()
	entries := []*tilepb.MapEntry{}
	for i := 0; i < 20; i++ {
		entries = append(entries, createEntry(fmt.Sprintf("key: %d", i), fmt.Sprintf("value: %d", i)))
	}

	tiles, err := DeNovo(s, beam.CreateList(s, entries), treeID, hash, 1)
	if err != nil {
		t.Fatalf("failed to create pipeline: %v", err)
	}

	passert.True(s, tiles, func(t *tilepb.MapTile) bool { return isStrictlySorted(t.Leaves) })

	if err := ptest.Run(p); err != nil {
		t.Fatalf("pipeline failed: %v", err)
	}
}

func TestGolden(t *testing.T) {
	p, s := beam.NewPipelineWithRoot()
	leaves := beam.CreateList(s, leafNodes(t, 500))

	tiles, err := DeNovo(s, leaves, 42, crypto.SHA256, 1)
	if err != nil {
		t.Fatalf("failed to create pipeline: %v", err)
	}
	rootTile := filter.Include(s, tiles, func(t *tilepb.MapTile) bool { return len(t.GetIndex()) == 0 })
	roots := beam.ParDo(s, func(t *tilepb.MapTile) string { return fmt.Sprintf("%x", t.Root) }, rootTile)

	passert.Equals(s, roots, "daf17dc2c83f37962bae8a65d294ef7fca4ffa02c10bdc4ca5c4dec408001c98")
	if err := ptest.Run(p); err != nil {
		t.Fatalf("pipeline failed: %v", err)
	}
}

// Ripped off from http://google3/third_party/golang/trillian/merkle/smt/hstar3_test.go?l=201&rcl=298994396
func leafNodes(t testing.TB, n int) []*tilepb.MapEntry {
	t.Helper()
	// Use a random sequence that depends on n.
	r := rand.New(rand.NewSource(int64(n)))
	entries := make([]*tilepb.MapEntry, n)
	for i := range entries {
		value := make([]byte, 32)
		if _, err := r.Read(value); err != nil {
			t.Fatalf("Failed to make random leaf hash: %v", err)
		}
		path := make([]byte, 32)
		if _, err := r.Read(path); err != nil {
			t.Fatalf("Failed to make random path: %v", err)
		}
		entries[i] = &tilepb.MapEntry{
			Key:   path,
			Value: value,
		}
	}

	return entries
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		vs    entriesIter
		empty bool
	}{
		{vs: entriesIter{}, empty: true},
		{vs: entriesIter{"a"}, empty: false},
		{vs: entriesIter{"a", "b", "c"}, empty: false},
	}

	for _, test := range tests {
		vs := test.vs
		empty, entries := isEmpty(vs.toIter())
		if got, want := empty, test.empty; got != want {
			t.Errorf("Emptiness mismatch, got != want (%t, %t)", got, want)
		}

		var gotKeys []string
		var entry *tilepb.MapEntry
		for entries(&entry) {
			gotKeys = append(gotKeys, string(entry.Key))
		}
		if got, want := len(gotKeys), len(vs); got != want {
			t.Errorf("got != want (%d, %d)", got, want)
		}
		for i := 0; i < len(vs); i++ {
			if got, want := gotKeys[i], vs[i]; got != want {
				t.Errorf("Mismatch at %d. got != want (%s, %s)", i, got, want)
			}
		}
	}
}

type entriesIter []string

func (i entriesIter) toIter() func(**tilepb.MapEntry) bool {
	pos := 0
	return func(r **tilepb.MapEntry) bool {
		if pos < len(i) {
			*r = &tilepb.MapEntry{
				Key:   []byte(i[pos]),
				Value: []byte(i[pos]),
			}
			pos++
			return true
		}
		return false
	}
}

func createEntry(k, v string) *tilepb.MapEntry {
	h := crypto.SHA256.New()
	h.Write([]byte(k))
	hk := h.Sum(nil)

	h = crypto.SHA256.New()
	h.Write([]byte(v))
	hv := h.Sum(nil)

	return &tilepb.MapEntry{
		Key:   hk,
		Value: hv,
	}
}

func isStrictlySorted(leaves []*tilepb.TileLeaf) bool {
	for i := 1; i < len(leaves); i++ {
		lPath, rPath := leaves[i-1].Path, leaves[i].Path
		if string(lPath) >= string(rPath) {
			return false
		}
	}
	return true
}
