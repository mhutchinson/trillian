package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/trillian"
	"github.com/google/trillian/monitoring"
	"github.com/google/trillian/monitoring/opencensus"
	"github.com/google/trillian/storage"
	"github.com/google/trillian/trees"
	"github.com/google/trillian/util"
	"k8s.io/klog/v2"
)

var (
	storageSystem = flag.String("storage_system", "mysql", fmt.Sprintf("Storage system to use. One of: %v", storage.Providers()))
	treeID        = flag.Int64("tree_id", 0, "The ID of the tree to extract tiles for.")
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	defer klog.Flush()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go util.AwaitSignal(ctx, cancel)

	mf := monitoring.InertMetricFactory{}
	monitoring.SetStartSpan(opencensus.StartSpan)

	sp, err := storage.NewProvider(*storageSystem, mf)
	if err != nil {
		klog.Exitf("Failed to get storage provider: %v", err)
	}
	defer sp.Close()

	optsLogRead := trees.NewGetOpts(trees.Query, trillian.TreeType_LOG, trillian.TreeType_PREORDERED_LOG)
	tree, err := trees.GetTree(ctx, sp.AdminStorage(), *treeID, optsLogRead)
	if err != nil {
		klog.Exitf("Failed to get tree %d: %v", *treeID, err)
	}

	tx, err := sp.LogStorage().SnapshotForTree(ctx, tree)
	if err != nil {
		klog.Exitf("Failed to get tree transaction: %v", err)
	}
	ttx, ok := tx.(storage.TileReadTx)
	if !ok {
		klog.Exitf("Storage layer doesn't support tile reading")
	}
	tileIDs := [][]byte{}
	ttx.GetSubtrees(ctx, tileIDs)
}
