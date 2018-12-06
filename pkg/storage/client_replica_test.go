// Copyright 2015 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package storage_test

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/base"
	"github.com/cockroachdb/cockroach/pkg/config"
	"github.com/cockroachdb/cockroach/pkg/internal/client"
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/server"
	"github.com/cockroachdb/cockroach/pkg/storage"
	"github.com/cockroachdb/cockroach/pkg/storage/batcheval"
	"github.com/cockroachdb/cockroach/pkg/storage/engine"
	"github.com/cockroachdb/cockroach/pkg/storage/engine/enginepb"
	"github.com/cockroachdb/cockroach/pkg/storage/storagebase"
	"github.com/cockroachdb/cockroach/pkg/storage/storagepb"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/serverutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/sqlutils"
	"github.com/cockroachdb/cockroach/pkg/testutils/testcluster"
	"github.com/cockroachdb/cockroach/pkg/util/caller"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/randutil"
	"github.com/cockroachdb/cockroach/pkg/util/stop"
	"github.com/cockroachdb/cockroach/pkg/util/syncutil"
	"github.com/pkg/errors"
)

// TestRangeCommandClockUpdate verifies that followers update their
// clocks when executing a command, even if the lease holder's clock is far
// in the future.
func TestRangeCommandClockUpdate(t *testing.T) {
	defer leaktest.AfterTest(t)()

	const numNodes = 3
	var manuals []*hlc.ManualClock
	var clocks []*hlc.Clock
	for i := 0; i < numNodes; i++ {
		manuals = append(manuals, hlc.NewManualClock(1))
		clocks = append(clocks, hlc.NewClock(manuals[i].UnixNano, 100*time.Millisecond))
	}
	mtc := &multiTestContext{clocks: clocks}
	defer mtc.Stop()
	mtc.Start(t, numNodes)
	mtc.replicateRange(1, 1, 2)

	// Advance the lease holder's clock ahead of the followers (by more than
	// MaxOffset but less than the range lease) and execute a command.
	manuals[0].Increment(int64(500 * time.Millisecond))
	incArgs := incrementArgs([]byte("a"), 5)
	ts := clocks[0].Now()
	if _, err := client.SendWrappedWith(context.Background(), mtc.stores[0].TestSender(), roachpb.Header{Timestamp: ts}, incArgs); err != nil {
		t.Fatal(err)
	}

	// Wait for that command to execute on all the followers.
	testutils.SucceedsSoon(t, func() error {
		values := []int64{}
		for _, eng := range mtc.engines {
			val, _, err := engine.MVCCGet(context.Background(), eng, roachpb.Key("a"), clocks[0].Now(), true, nil)
			if err != nil {
				return err
			}
			values = append(values, mustGetInt(val))
		}
		if !reflect.DeepEqual(values, []int64{5, 5, 5}) {
			return errors.Errorf("expected (5, 5, 5), got %v", values)
		}
		return nil
	})

	// Verify that all the followers have accepted the clock update from
	// node 0 even though it comes from outside the usual max offset.
	now := clocks[0].Now()
	for i, clock := range clocks {
		// Only compare the WallTimes: it's normal for clock 0 to be a few logical ticks ahead.
		if clock.Now().WallTime < now.WallTime {
			t.Errorf("clock %d is behind clock 0: %s vs %s", i, clock.Now(), now)
		}
	}
}

// TestRejectFutureCommand verifies that lease holders reject commands that
// would cause a large time jump.
func TestRejectFutureCommand(t *testing.T) {
	defer leaktest.AfterTest(t)()

	manual := hlc.NewManualClock(123)
	clock := hlc.NewClock(manual.UnixNano, 100*time.Millisecond)
	mtc := &multiTestContext{clock: clock}
	defer mtc.Stop()
	mtc.Start(t, 1)

	ts1 := clock.Now()

	key := roachpb.Key("a")
	incArgs := incrementArgs(key, 5)

	// Commands with a future timestamp that is within the MaxOffset
	// bound will be accepted and will cause the clock to advance.
	const numCmds = 3
	clockOffset := clock.MaxOffset() / numCmds
	for i := int64(1); i <= numCmds; i++ {
		ts := ts1.Add(i*clockOffset.Nanoseconds(), 0)
		if _, err := client.SendWrappedWith(context.Background(), mtc.stores[0].TestSender(), roachpb.Header{Timestamp: ts}, incArgs); err != nil {
			t.Fatal(err)
		}
	}

	ts2 := clock.Now()
	if expAdvance, advance := ts2.GoTime().Sub(ts1.GoTime()), numCmds*clockOffset; advance != expAdvance {
		t.Fatalf("expected clock to advance %s; got %s", expAdvance, advance)
	}

	// Once the accumulated offset reaches MaxOffset, commands will be rejected.
	_, pErr := client.SendWrappedWith(context.Background(), mtc.stores[0].TestSender(), roachpb.Header{Timestamp: ts1.Add(clock.MaxOffset().Nanoseconds()+1, 0)}, incArgs)
	if !testutils.IsPError(pErr, "remote wall time is too far ahead") {
		t.Fatalf("unexpected error %v", pErr)
	}

	// The clock did not advance and the final command was not executed.
	ts3 := clock.Now()
	if advance := ts3.GoTime().Sub(ts2.GoTime()); advance != 0 {
		t.Fatalf("expected clock not to advance, but it advanced by %s", advance)
	}
	val, _, err := engine.MVCCGet(context.Background(), mtc.engines[0], key, ts3, true, nil)
	if err != nil {
		t.Fatal(err)
	}
	if a, e := mustGetInt(val), incArgs.Increment*numCmds; a != e {
		t.Errorf("expected %d, got %d", e, a)
	}
}

// TestTxnPutOutOfOrder tests a case where a put operation of an older
// timestamp comes after a put operation of a newer timestamp in a
// txn. The test ensures such an out-of-order put succeeds and
// overrides an old value. The test uses a "Writer" and a "Reader"
// to reproduce an out-of-order put.
//
// 1) The Writer executes a cput operation and writes a write intent with
//    time T in a txn.
// 2) Before the Writer's txn is committed, the Reader sends a high priority
//    get operation with time T+100. This pushes the Writer txn timestamp to
//    T+100. The Reader also writes to the same key the Writer did a cput to
//    in order to trigger the restart of the Writer's txn. The original
//    write intent timestamp is also updated to T+100.
// 3) The Writer starts a new epoch of the txn, but before it writes, the
//    Reader sends another high priority get operation with time T+200. This
//    pushes the Writer txn timestamp to T+200 to trigger a restart of the
//    Writer txn. The Writer will not actually restart until it tries to commit
//    the current epoch of the transaction. The Reader updates the timestamp of
//    the write intent to T+200. The test deliberately fails the Reader get
//    operation, and cockroach doesn't update its read timestamp cache.
// 4) The Writer executes the put operation again. This put operation comes
//    out-of-order since its timestamp is T+100, while the intent timestamp
//    updated at Step 3 is T+200.
// 5) The put operation overrides the old value using timestamp T+100.
// 6) When the Writer attempts to commit its txn, the txn will be restarted
//    again at a new epoch timestamp T+200, which will finally succeed.
func TestTxnPutOutOfOrder(t *testing.T) {
	defer leaktest.AfterTest(t)()

	// key is selected to fall within the meta range in order for the later
	// routing of requests to range 1 to work properly. Removing the routing
	// of all requests to range 1 would allow us to make the key more normal.
	const (
		key        = "key"
		restartKey = "restart"
	)
	// Set up a filter to so that the get operation at Step 3 will return an error.
	var numGets int32

	stopper := stop.NewStopper()
	defer stopper.Stop(context.TODO())
	manual := hlc.NewManualClock(123)
	cfg := storage.TestStoreConfig(hlc.NewClock(manual.UnixNano, time.Nanosecond))
	// Splits can cause our chosen key to end up on a range other than range 1,
	// and trying to handle that complicates the test without providing any
	// added benefit.
	cfg.TestingKnobs.DisableSplitQueue = true
	cfg.TestingKnobs.EvalKnobs.TestingEvalFilter =
		func(filterArgs storagebase.FilterArgs) *roachpb.Error {
			if _, ok := filterArgs.Req.(*roachpb.GetRequest); ok &&
				filterArgs.Req.Header().Key.Equal(roachpb.Key(key)) &&
				filterArgs.Hdr.Txn == nil {
				// The Reader executes two get operations, each of which triggers two get requests
				// (the first request fails and triggers txn push, and then the second request
				// succeeds). Returns an error for the fourth get request to avoid timestamp cache
				// update after the third get operation pushes the txn timestamp.
				if atomic.AddInt32(&numGets, 1) == 4 {
					return roachpb.NewErrorWithTxn(errors.Errorf("Test"), filterArgs.Hdr.Txn)
				}
			}
			return nil
		}
	eng := engine.NewInMem(roachpb.Attributes{}, 10<<20)
	stopper.AddCloser(eng)
	store := createTestStoreWithEngine(t,
		eng,
		true,
		cfg,
		stopper,
	)

	// Put an initial value.
	initVal := []byte("initVal")
	err := store.DB().Put(context.TODO(), key, initVal)
	if err != nil {
		t.Fatalf("failed to put: %s", err)
	}

	waitPut := make(chan struct{})
	waitFirstGet := make(chan struct{})
	waitTxnRestart := make(chan struct{})
	waitSecondGet := make(chan struct{})
	errChan := make(chan error)

	// Start the Writer.
	go func() {
		epoch := -1
		// Start a txn that does read-after-write.
		// The txn will be restarted twice, and the out-of-order put
		// will happen in the second epoch.
		errChan <- store.DB().Txn(context.TODO(), func(ctx context.Context, txn *client.Txn) error {
			epoch++

			if epoch == 1 {
				// Wait until the second get operation is issued.
				close(waitTxnRestart)
				<-waitSecondGet
			}

			// Get a key which we can write to from the Reader in order to force a restart.
			if _, err := txn.Get(ctx, restartKey); err != nil {
				return err
			}

			updatedVal := []byte("updatedVal")
			if err := txn.CPut(ctx, key, updatedVal, "initVal"); err != nil {
				log.Errorf(context.TODO(), "failed put value: %s", err)
				return err
			}

			// Make sure a get will return the value that was just written.
			actual, err := txn.Get(ctx, key)
			if err != nil {
				return err
			}
			if !bytes.Equal(actual.ValueBytes(), updatedVal) {
				return errors.Errorf("unexpected get result: %s", actual)
			}

			if epoch == 0 {
				// Wait until the first get operation will push the txn timestamp.
				close(waitPut)
				<-waitFirstGet
			}

			b := txn.NewBatch()
			return txn.CommitInBatch(ctx, b)
		})

		if epoch != 2 {
			file, line, _ := caller.Lookup(0)
			errChan <- errors.Errorf("%s:%d unexpected number of txn retries. "+
				"Expected epoch 2, got: %d.", file, line, epoch)
		} else {
			errChan <- nil
		}
	}()

	<-waitPut

	// Start the Reader.

	// Advance the clock and send a get operation with higher
	// priority to trigger the txn restart.
	manual.Increment(100)

	priority := roachpb.UserPriority(-math.MaxInt32)
	requestHeader := roachpb.RequestHeader{
		Key: roachpb.Key(key),
	}
	h := roachpb.Header{
		Timestamp:    cfg.Clock.Now(),
		UserPriority: priority,
	}
	if _, err := client.SendWrappedWith(
		context.Background(), store.TestSender(), h, &roachpb.GetRequest{RequestHeader: requestHeader},
	); err != nil {
		t.Fatalf("failed to get: %s", err)
	}
	// Write to the restart key so that the Writer's txn must restart.
	putReq := &roachpb.PutRequest{
		RequestHeader: roachpb.RequestHeader{Key: roachpb.Key(restartKey)},
		Value:         roachpb.MakeValueFromBytes([]byte("restart-value")),
	}
	if _, err := client.SendWrappedWith(context.Background(), store.TestSender(), h, putReq); err != nil {
		t.Fatalf("failed to put: %s", err)
	}

	// Wait until the writer restarts the txn.
	close(waitFirstGet)
	<-waitTxnRestart

	// Advance the clock and send a get operation again. This time
	// we use TestingCommandFilter so that a get operation is not
	// processed after the write intent is resolved (to prevent the
	// timestamp cache from being updated).
	manual.Increment(100)

	h.Timestamp = cfg.Clock.Now()
	if _, err := client.SendWrappedWith(
		context.Background(), store.TestSender(), h, &roachpb.GetRequest{RequestHeader: requestHeader},
	); err == nil {
		t.Fatal("unexpected success of get")
	}
	if _, err := client.SendWrappedWith(context.Background(), store.TestSender(), h, putReq); err != nil {
		t.Fatalf("failed to put: %s", err)
	}

	close(waitSecondGet)
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			t.Fatal(err)
		}
	}
}

// TestRangeLookupUseReverse tests whether the results and the results count
// are correct when scanning in reverse order.
func TestRangeLookupUseReverse(t *testing.T) {
	defer leaktest.AfterTest(t)()
	storeCfg := storage.TestStoreConfig(nil)
	storeCfg.TestingKnobs.DisableSplitQueue = true
	storeCfg.TestingKnobs.DisableMergeQueue = true
	stopper := stop.NewStopper()
	defer stopper.Stop(context.TODO())
	store := createTestStoreWithConfig(t, stopper, storeCfg)

	// Init test ranges:
	// ["","a"), ["a","c"), ["c","e"), ["e","g") and ["g","\xff\xff").
	splits := []*roachpb.AdminSplitRequest{
		adminSplitArgs(roachpb.Key("g")),
		adminSplitArgs(roachpb.Key("e")),
		adminSplitArgs(roachpb.Key("c")),
		adminSplitArgs(roachpb.Key("a")),
	}

	for _, split := range splits {
		_, pErr := client.SendWrapped(context.Background(), store.TestSender(), split)
		if pErr != nil {
			t.Fatalf("%q: split unexpected error: %s", split.SplitKey, pErr)
		}
	}

	// Resolve the intents.
	scanArgs := roachpb.ScanRequest{
		RequestHeader: roachpb.RequestHeader{
			Key:    keys.RangeMetaKey(roachpb.RKeyMin.Next()).AsRawKey(),
			EndKey: keys.RangeMetaKey(roachpb.RKeyMax).AsRawKey(),
		},
	}
	testutils.SucceedsSoon(t, func() error {
		_, pErr := client.SendWrapped(context.Background(), store.TestSender(), &scanArgs)
		return pErr.GoError()
	})

	testCases := []struct {
		key         roachpb.RKey
		maxResults  int64
		expected    []roachpb.RangeDescriptor
		expectedPre []roachpb.RangeDescriptor
	}{
		// Test key in the middle of the range.
		{
			key:        roachpb.RKey("f"),
			maxResults: 2,
			// ["e","g") and ["c","e").
			expected: []roachpb.RangeDescriptor{
				{StartKey: roachpb.RKey("e"), EndKey: roachpb.RKey("g")},
			},
			expectedPre: []roachpb.RangeDescriptor{
				{StartKey: roachpb.RKey("c"), EndKey: roachpb.RKey("e")},
			},
		},
		// Test key in the end key of the range.
		{
			key:        roachpb.RKey("g"),
			maxResults: 3,
			// ["e","g"), ["c","e") and ["a","c").
			expected: []roachpb.RangeDescriptor{
				{StartKey: roachpb.RKey("e"), EndKey: roachpb.RKey("g")},
			},
			expectedPre: []roachpb.RangeDescriptor{
				{StartKey: roachpb.RKey("c"), EndKey: roachpb.RKey("e")},
				{StartKey: roachpb.RKey("a"), EndKey: roachpb.RKey("c")},
			},
		},
		{
			key:        roachpb.RKey("e"),
			maxResults: 2,
			// ["c","e") and ["a","c").
			expected: []roachpb.RangeDescriptor{
				{StartKey: roachpb.RKey("c"), EndKey: roachpb.RKey("e")},
			},
			expectedPre: []roachpb.RangeDescriptor{
				{StartKey: roachpb.RKey("a"), EndKey: roachpb.RKey("c")},
			},
		},
		// Test RKeyMax.
		{
			key:        roachpb.RKeyMax,
			maxResults: 2,
			// ["e","g") and ["g","\xff\xff")
			expected: []roachpb.RangeDescriptor{
				{StartKey: roachpb.RKey("g"), EndKey: roachpb.RKey("\xff\xff")},
			},
			expectedPre: []roachpb.RangeDescriptor{
				{StartKey: roachpb.RKey("e"), EndKey: roachpb.RKey("g")},
			},
		},
		// Test Meta2KeyMax.
		{
			key:        roachpb.RKey(keys.Meta2KeyMax),
			maxResults: 1,
			// ["","a")
			expected: []roachpb.RangeDescriptor{
				{StartKey: roachpb.RKeyMin, EndKey: roachpb.RKey("a")},
			},
		},
	}

	for _, test := range testCases {
		t.Run(fmt.Sprintf("key=%s", test.key), func(t *testing.T) {
			rs, preRs, err := client.RangeLookup(context.Background(), store.TestSender(),
				test.key.AsRawKey(), roachpb.READ_UNCOMMITTED, test.maxResults-1, true /* prefetchReverse */)
			if err != nil {
				t.Fatalf("LookupRange error: %s", err)
			}

			// Checks the results count.
			if rsLen, preRsLen := len(rs), len(preRs); int64(rsLen+preRsLen) != test.maxResults {
				t.Fatalf("returned results count, expected %d, but got %d+%d", test.maxResults, rsLen, preRsLen)
			}
			// Checks the range descriptors.
			for _, rngSlice := range []struct {
				expect, reply []roachpb.RangeDescriptor
			}{
				{test.expected, rs},
				{test.expectedPre, preRs},
			} {
				for i, rng := range rngSlice.expect {
					if !(rng.StartKey.Equal(rngSlice.reply[i].StartKey) && rng.EndKey.Equal(rngSlice.reply[i].EndKey)) {
						t.Fatalf("returned range is not correct, expected %v, but got %v", rng, rngSlice.reply[i])
					}
				}
			}
		})
	}
}

type leaseTransferTest struct {
	mtc                        *multiTestContext
	replica0, replica1         *storage.Replica
	replica0Desc, replica1Desc roachpb.ReplicaDescriptor
	leftKey                    roachpb.Key
	filterMu                   syncutil.Mutex
	filter                     func(filterArgs storagebase.FilterArgs) *roachpb.Error
	waitForTransferBlocked     atomic.Value
	transferBlocked            chan struct{}
}

func setupLeaseTransferTest(t *testing.T) *leaseTransferTest {
	l := &leaseTransferTest{
		leftKey: roachpb.Key("a"),
	}

	cfg := storage.TestStoreConfig(nil)
	// Ensure the node liveness duration isn't too short. By default it is 900ms
	// for TestStoreConfig().
	cfg.RangeLeaseRaftElectionTimeoutMultiplier =
		float64((9 * time.Second) / cfg.RaftElectionTimeout())
	cfg.TestingKnobs.EvalKnobs.TestingEvalFilter =
		func(filterArgs storagebase.FilterArgs) *roachpb.Error {
			l.filterMu.Lock()
			filterCopy := l.filter
			l.filterMu.Unlock()
			if filterCopy != nil {
				return filterCopy(filterArgs)
			}
			return nil
		}

	l.waitForTransferBlocked.Store(false)
	l.transferBlocked = make(chan struct{})
	cfg.TestingKnobs.LeaseTransferBlockedOnExtensionEvent = func(
		_ roachpb.ReplicaDescriptor) {
		if l.waitForTransferBlocked.Load().(bool) {
			l.transferBlocked <- struct{}{}
			l.waitForTransferBlocked.Store(false)
		}
	}

	l.mtc = &multiTestContext{}
	l.mtc.storeConfig = &cfg
	l.mtc.Start(t, 2)
	l.mtc.initGossipNetwork()

	// First, do a write; we'll use it to determine when the dust has settled.
	l.leftKey = roachpb.Key("a")
	incArgs := incrementArgs(l.leftKey, 1)
	if _, pErr := client.SendWrapped(context.Background(), l.mtc.distSenders[0], incArgs); pErr != nil {
		t.Fatal(pErr)
	}

	// Get the left range's ID.
	rangeID := l.mtc.stores[0].LookupReplica(keys.MustAddr(l.leftKey)).RangeID

	// Replicate the left range onto node 1.
	l.mtc.replicateRange(rangeID, 1)

	l.replica0 = l.mtc.stores[0].LookupReplica(roachpb.RKey("a"))
	l.replica1 = l.mtc.stores[1].LookupReplica(roachpb.RKey("a"))
	{
		var err error
		if l.replica0Desc, err = l.replica0.GetReplicaDescriptor(); err != nil {
			t.Fatal(err)
		}
		if l.replica1Desc, err = l.replica1.GetReplicaDescriptor(); err != nil {
			t.Fatal(err)
		}
	}

	// Check that replica0 can serve reads OK.
	if pErr := l.sendRead(0); pErr != nil {
		t.Fatal(pErr)
	}
	return l
}

func (l *leaseTransferTest) sendRead(storeIdx int) *roachpb.Error {
	desc := l.mtc.stores[storeIdx].LookupReplica(keys.MustAddr(l.leftKey))
	replicaDesc, err := desc.GetReplicaDescriptor()
	if err != nil {
		return roachpb.NewError(err)
	}
	_, pErr := client.SendWrappedWith(
		context.Background(),
		l.mtc.senders[storeIdx],
		roachpb.Header{RangeID: desc.RangeID, Replica: replicaDesc},
		getArgs(l.leftKey),
	)
	if pErr != nil {
		log.Warning(context.TODO(), pErr)
	}
	return pErr
}

// checkHasLease checks that a lease for the left range is owned by a
// replica. The check is executed in a retry loop because the lease may not
// have been applied yet.
func (l *leaseTransferTest) checkHasLease(t *testing.T, storeIdx int) {
	t.Helper()
	testutils.SucceedsSoon(t, func() error {
		return l.sendRead(storeIdx).GoError()
	})
}

// setFilter is a helper function to enable/disable the blocking of
// RequestLeaseRequests on replica1. This function will notify that an
// extension is blocked on the passed in channel and will wait on the same
// channel to unblock the extension. Note that once an extension is blocked,
// the filter is cleared.
func (l *leaseTransferTest) setFilter(setTo bool, extensionSem chan struct{}) {
	l.filterMu.Lock()
	defer l.filterMu.Unlock()
	if !setTo {
		l.filter = nil
		return
	}
	l.filter = func(filterArgs storagebase.FilterArgs) *roachpb.Error {
		if filterArgs.Sid != l.mtc.stores[1].Ident.StoreID {
			return nil
		}
		llReq, ok := filterArgs.Req.(*roachpb.RequestLeaseRequest)
		if !ok {
			return nil
		}
		if llReq.Lease.Replica == l.replica1Desc {
			// Notify the main thread that the extension is in progress and wait for
			// the signal to proceed.
			l.filterMu.Lock()
			l.filter = nil
			l.filterMu.Unlock()
			extensionSem <- struct{}{}
			<-extensionSem
		}
		return nil
	}
}

func (l *leaseTransferTest) forceLeaseExtension(storeIdx int, lease roachpb.Lease) error {
	shouldRenewTS := lease.Expiration.Add(-1, 0)
	l.mtc.manualClock.Set(shouldRenewTS.WallTime + 1)
	err := l.sendRead(storeIdx).GoError()
	if err != nil {
		// We can sometimes receive an error from our renewal attempt because the
		// lease transfer ends up causing the renewal to re-propose and second
		// attempt fails because it's already been renewed. This used to work
		// before we compared the proposer's lease with the actual lease because
		// the renewed lease still encompassed the previous request.
		if _, ok := err.(*roachpb.NotLeaseHolderError); ok {
			err = nil
		}
	}
	return err
}

// ensureLeaderAndRaftState is a helper function that blocks until leader is
// the raft leader and follower is up to date.
func (l *leaseTransferTest) ensureLeaderAndRaftState(
	t *testing.T, leader *storage.Replica, follower roachpb.ReplicaDescriptor,
) {
	t.Helper()
	leaderDesc, err := leader.GetReplicaDescriptor()
	if err != nil {
		t.Fatal(err)
	}
	testutils.SucceedsSoon(t, func() error {
		r := l.mtc.getRaftLeader(l.replica0.RangeID)
		if r == nil {
			return errors.Errorf("could not find raft leader replica for range %d", l.replica0.RangeID)
		}
		desc, err := r.GetReplicaDescriptor()
		if err != nil {
			return errors.Wrap(err, "could not get replica descriptor")
		}
		if desc != leaderDesc {
			return errors.Errorf(
				"expected replica with id %v to be raft leader, instead got id %v",
				leaderDesc.ReplicaID,
				desc.ReplicaID,
			)
		}
		return nil
	})

	testutils.SucceedsSoon(t, func() error {
		status := leader.RaftStatus()
		progress, ok := status.Progress[uint64(follower.ReplicaID)]
		if !ok {
			return errors.Errorf(
				"replica %v progress not found in progress map: %v",
				follower.ReplicaID,
				status.Progress,
			)
		}
		if progress.Match < status.Commit {
			return errors.Errorf("replica %v failed to catch up", follower.ReplicaID)
		}
		return nil
	})
}

func TestRangeTransferLeaseExpirationBased(t *testing.T) {
	defer leaktest.AfterTest(t)()

	t.Run("Transfer", func(t *testing.T) {
		l := setupLeaseTransferTest(t)
		defer l.mtc.Stop()
		origLease, _ := l.replica0.GetLease()
		{
			// Transferring the lease to ourself should be a no-op.
			if err := l.replica0.AdminTransferLease(context.Background(), l.replica0Desc.StoreID); err != nil {
				t.Fatal(err)
			}
			newLease, _ := l.replica0.GetLease()
			if !origLease.Equivalent(newLease) {
				t.Fatalf("original lease %v and new lease %v not equivalent", origLease, newLease)
			}
		}

		{
			// An invalid target should result in an error.
			const expected = "unable to find store .* in range"
			if err := l.replica0.AdminTransferLease(context.Background(), 1000); !testutils.IsError(err, expected) {
				t.Fatalf("expected %s, but found %v", expected, err)
			}
		}

		if err := l.replica0.AdminTransferLease(context.Background(), l.replica1Desc.StoreID); err != nil {
			t.Fatal(err)
		}

		// Check that replica0 doesn't serve reads any more.
		pErr := l.sendRead(0)
		nlhe, ok := pErr.GetDetail().(*roachpb.NotLeaseHolderError)
		if !ok {
			t.Fatalf("expected %T, got %s", &roachpb.NotLeaseHolderError{}, pErr)
		}
		if *(nlhe.LeaseHolder) != l.replica1Desc {
			t.Fatalf("expected lease holder %+v, got %+v",
				l.replica1Desc, nlhe.LeaseHolder)
		}

		// Check that replica1 now has the lease.
		l.checkHasLease(t, 1)

		replica1Lease, _ := l.replica1.GetLease()

		// We'd like to verify the timestamp cache's low water mark, but this is
		// impossible to determine precisely in all cases because it may have
		// been subsumed by future tscache accesses. So instead of checking the
		// low water mark, we make sure that the high water mark is equal to or
		// greater than the new lease start time, which is less than the
		// previous lease's expiration time.
		if highWater := l.replica1.GetTSCacheHighWater(); highWater.Less(replica1Lease.Start) {
			t.Fatalf("expected timestamp cache high water %s, but found %s",
				replica1Lease.Start, highWater)
		}
	})

	// Make replica1 extend its lease and transfer the lease immediately after
	// that. Test that the transfer still happens (it'll wait until the extension
	// is done).
	t.Run("TransferWithExtension", func(t *testing.T) {
		l := setupLeaseTransferTest(t)
		defer l.mtc.Stop()
		// Ensure that replica1 has the lease.
		if err := l.replica0.AdminTransferLease(context.Background(), l.replica1Desc.StoreID); err != nil {
			t.Fatal(err)
		}
		l.checkHasLease(t, 1)

		extensionSem := make(chan struct{})
		l.setFilter(true, extensionSem)

		// Initiate an extension.
		renewalErrCh := make(chan error)
		go func() {
			lease, _ := l.replica1.GetLease()
			renewalErrCh <- l.forceLeaseExtension(1, lease)
		}()

		// Wait for extension to be blocked.
		<-extensionSem
		l.waitForTransferBlocked.Store(true)
		// Initiate a transfer.
		transferErrCh := make(chan error)
		go func() {
			// Transfer back from replica1 to replica0.
			err := l.replica1.AdminTransferLease(context.Background(), l.replica0Desc.StoreID)
			// Ignore not leaseholder errors which can arise due to re-proposals.
			if _, ok := err.(*roachpb.NotLeaseHolderError); ok {
				err = nil
			}
			transferErrCh <- err
		}()
		// Wait for the transfer to be blocked by the extension.
		<-l.transferBlocked
		// Now unblock the extension.
		extensionSem <- struct{}{}
		l.checkHasLease(t, 0)
		l.setFilter(false, nil)

		if err := <-renewalErrCh; err != nil {
			t.Errorf("unexpected error from lease renewal: %s", err)
		}
		if err := <-transferErrCh; err != nil {
			t.Errorf("unexpected error from lease transfer: %s", err)
		}
	})

	// DrainTransfer verifies that a draining store attempts to transfer away
	// range leases owned by its replicas.
	t.Run("DrainTransfer", func(t *testing.T) {
		l := setupLeaseTransferTest(t)
		defer l.mtc.Stop()
		// We have to ensure that replica0 is the raft leader and that replica1 has
		// caught up to replica0 as draining code doesn't transfer leases to
		// behind replicas.
		l.ensureLeaderAndRaftState(t, l.replica0, l.replica1Desc)
		l.mtc.stores[0].SetDraining(true)

		// Check that replica0 doesn't serve reads any more.
		pErr := l.sendRead(0)
		nlhe, ok := pErr.GetDetail().(*roachpb.NotLeaseHolderError)
		if !ok {
			t.Fatalf("expected %T, got %s", &roachpb.NotLeaseHolderError{}, pErr)
		}
		if nlhe.LeaseHolder == nil || *nlhe.LeaseHolder != l.replica1Desc {
			t.Fatalf("expected lease holder %+v, got %+v",
				l.replica1Desc, nlhe.LeaseHolder)
		}

		// Check that replica1 now has the lease.
		l.checkHasLease(t, 1)

		l.mtc.stores[0].SetDraining(false)
	})

	// DrainTransferWithExtension verifies that a draining store waits for any
	// in-progress lease requests to complete before transferring away the new
	// lease.
	t.Run("DrainTransferWithExtension", func(t *testing.T) {
		l := setupLeaseTransferTest(t)
		defer l.mtc.Stop()
		// Ensure that replica1 has the lease.
		if err := l.replica0.AdminTransferLease(context.Background(), l.replica1Desc.StoreID); err != nil {
			t.Fatal(err)
		}
		l.checkHasLease(t, 1)

		extensionSem := make(chan struct{})
		l.setFilter(true, extensionSem)

		// Initiate an extension.
		renewalErrCh := make(chan error)
		go func() {
			lease, _ := l.replica1.GetLease()
			renewalErrCh <- l.forceLeaseExtension(1, lease)
		}()

		// Wait for extension to be blocked.
		<-extensionSem

		// Make sure that replica 0 is up to date enough to receive the lease.
		l.ensureLeaderAndRaftState(t, l.replica1, l.replica0Desc)

		// Drain node 1 with an extension in progress.
		go func() {
			l.mtc.stores[1].SetDraining(true)
		}()
		// Now unblock the extension.
		extensionSem <- struct{}{}

		l.checkHasLease(t, 0)
		l.setFilter(false, nil)

		if err := <-renewalErrCh; err != nil {
			t.Errorf("unexpected error from lease renewal: %s", err)
		}
	})
}

// TestRangeLimitTxnMaxTimestamp verifies that on lease transfer, the
// normal limiting of a txn's max timestamp to the first observed
// timestamp on a node is extended to include the lease start
// timestamp. This disallows the possibility that a write to another
// replica of the range (on node n1) happened at a later timestamp
// than the originally observed timestamp for the node which now owns
// the lease (n2). This can happen if the replication of the write
// doesn't make it from n1 to n2 before the transaction observes n2's
// clock time.
func TestRangeLimitTxnMaxTimestamp(t *testing.T) {
	defer leaktest.AfterTest(t)()
	cfg := storage.TestStoreConfig(nil)
	cfg.RangeLeaseRaftElectionTimeoutMultiplier =
		float64((9 * time.Second) / cfg.RaftElectionTimeout())
	mtc := &multiTestContext{}
	mtc.storeConfig = &cfg
	keyA := roachpb.Key("a")
	// Create a new clock for node2 to allow drift between the two wall clocks.
	manual1 := hlc.NewManualClock(100) // node1 clock is @t=100
	clock1 := hlc.NewClock(manual1.UnixNano, 250*time.Nanosecond)
	manual2 := hlc.NewManualClock(98) // node2 clock is @t=98
	clock2 := hlc.NewClock(manual2.UnixNano, 250*time.Nanosecond)
	mtc.clocks = []*hlc.Clock{clock1, clock2}

	// Start a transaction using node2 as a gateway.
	txn := roachpb.MakeTransaction("test", keyA, 1, enginepb.SERIALIZABLE, clock2.Now(), 250 /* maxOffsetNs */)
	// Simulate a read to another range on node2 by setting the observed timestamp.
	txn.UpdateObservedTimestamp(2, clock2.Now())

	defer mtc.Stop()
	mtc.Start(t, 2)

	// Do a write on node1 to establish a key with its timestamp @t=100.
	if _, pErr := client.SendWrapped(
		context.Background(), mtc.distSenders[0], putArgs(keyA, []byte("value")),
	); pErr != nil {
		t.Fatal(pErr)
	}

	// Up-replicate the data in the range to node2.
	replica1 := mtc.stores[0].LookupReplica(roachpb.RKey(keyA))
	mtc.replicateRange(replica1.RangeID, 1)

	// Transfer the lease from node1 to node2.
	replica2 := mtc.stores[1].LookupReplica(roachpb.RKey(keyA))
	replica2Desc, err := replica2.GetReplicaDescriptor()
	if err != nil {
		t.Fatal(err)
	}
	testutils.SucceedsSoon(t, func() error {
		if err := replica1.AdminTransferLease(context.Background(), replica2Desc.StoreID); err != nil {
			t.Fatal(err)
		}
		lease, _ := replica2.GetLease()
		if lease.Replica.NodeID != replica2.NodeID() {
			return errors.Errorf("expected lease transfer to node2: %s", lease)
		}
		return nil
	})
	// Verify that after the lease transfer, node2's clock has advanced to at least 100.
	if now1, now2 := clock1.Now(), clock2.Now(); now2.WallTime < now1.WallTime {
		t.Fatalf("expected node2's clock walltime to be >= %d; got %d", now1.WallTime, now2.WallTime)
	}

	// Send a get request for keyA to node2, which is now the
	// leaseholder. If the max timestamp were not being properly limited,
	// we would end up incorrectly reading nothing for keyA. Instead we
	// expect to see an uncertainty interval error.
	h := roachpb.Header{Txn: &txn}
	if _, pErr := client.SendWrappedWith(
		context.Background(), mtc.distSenders[0], h, getArgs(keyA),
	); !testutils.IsPError(pErr, "uncertainty") {
		t.Fatalf("expected an uncertainty interval error; got %v", pErr)
	}
}

// TestLeaseMetricsOnSplitAndTransfer verifies that lease-related metrics
// are updated after splitting a range and then initiating one successful
// and one failing lease transfer.
func TestLeaseMetricsOnSplitAndTransfer(t *testing.T) {
	defer leaktest.AfterTest(t)()
	var injectLeaseTransferError atomic.Value
	sc := storage.TestStoreConfig(nil)
	sc.TestingKnobs.DisableSplitQueue = true
	sc.TestingKnobs.DisableMergeQueue = true
	sc.TestingKnobs.EvalKnobs.TestingEvalFilter =
		func(filterArgs storagebase.FilterArgs) *roachpb.Error {
			if args, ok := filterArgs.Req.(*roachpb.TransferLeaseRequest); ok {
				if val := injectLeaseTransferError.Load(); val != nil && val.(bool) {
					// Note that we can't just return an error here as we only
					// end up counting failures in the metrics if the command
					// makes it through to being executed. So use a fake store ID.
					args.Lease.Replica.StoreID = roachpb.StoreID(1000)
				}
			}
			return nil
		}
	mtc := &multiTestContext{storeConfig: &sc}
	defer mtc.Stop()
	mtc.Start(t, 2)

	// Up-replicate to two replicas.
	keyMinReplica0 := mtc.stores[0].LookupReplica(roachpb.RKeyMin)
	mtc.replicateRange(keyMinReplica0.RangeID, 1)

	// Split the key space at key "a".
	splitKey := roachpb.RKey("a")
	splitArgs := adminSplitArgs(splitKey.AsRawKey())
	if _, pErr := client.SendWrapped(
		context.Background(), mtc.stores[0].TestSender(), splitArgs,
	); pErr != nil {
		t.Fatal(pErr)
	}

	// Now, a successful transfer from LHS replica 0 to replica 1.
	injectLeaseTransferError.Store(false)
	if err := mtc.dbs[0].AdminTransferLease(
		context.TODO(), keyMinReplica0.Desc().StartKey.AsRawKey(), mtc.stores[1].StoreID(),
	); err != nil {
		t.Fatalf("unable to transfer lease to replica 1: %s", err)
	}
	// Wait for all replicas to process.
	testutils.SucceedsSoon(t, func() error {
		for i := 0; i < 2; i++ {
			r := mtc.stores[i].LookupReplica(roachpb.RKeyMin)
			if l, _ := r.GetLease(); l.Replica.StoreID != mtc.stores[1].StoreID() {
				return errors.Errorf("expected lease to transfer to replica 2: got %s", l)
			}
		}
		return nil
	})

	// Next a failed transfer from RHS replica 0 to replica 1.
	injectLeaseTransferError.Store(true)
	keyAReplica0 := mtc.stores[0].LookupReplica(splitKey)
	if err := mtc.dbs[0].AdminTransferLease(
		context.TODO(), keyAReplica0.Desc().StartKey.AsRawKey(), mtc.stores[1].StoreID(),
	); err == nil {
		t.Fatal("expected an error transferring to an unknown store ID")
	}

	metrics := mtc.stores[0].Metrics()
	if a, e := metrics.LeaseTransferSuccessCount.Count(), int64(1); a != e {
		t.Errorf("expected %d lease transfer successes; got %d", e, a)
	}
	if a, e := metrics.LeaseTransferErrorCount.Count(), int64(1); a != e {
		t.Errorf("expected %d lease transfer errors; got %d", e, a)
	}

	// Expire current leases and put a key to RHS of split to request
	// an epoch-based lease.
	testutils.SucceedsSoon(t, func() error {
		mtc.advanceClock(context.TODO())
		if err := mtc.stores[0].DB().Put(context.TODO(), "a", "foo"); err != nil {
			return err
		}

		// Update replication gauges for all stores and verify we have 1 each of
		// expiration and epoch leases.
		var expirationLeases int64
		var epochLeases int64
		for i := range mtc.stores {
			if err := mtc.stores[i].ComputeMetrics(context.Background(), 0); err != nil {
				return err
			}
			metrics = mtc.stores[i].Metrics()
			expirationLeases += metrics.LeaseExpirationCount.Value()
			epochLeases += metrics.LeaseEpochCount.Value()
		}
		if a, e := expirationLeases, int64(1); a != e {
			return errors.Errorf("expected %d expiration lease count; got %d", e, a)
		}
		if a, e := epochLeases, int64(1); a != e {
			return errors.Errorf("expected %d epoch lease count; got %d", e, a)
		}
		return nil
	})
}

// Test that leases held before a restart are not used after the restart.
// See replica.mu.minLeaseProposedTS for the reasons why this isn't allowed.
func TestLeaseNotUsedAfterRestart(t *testing.T) {
	defer leaktest.AfterTest(t)()
	ctx := context.Background()

	sc := storage.TestStoreConfig(nil)
	var leaseAcquisitionTrap atomic.Value
	// Disable the split queue so that no ranges are split. This makes it easy
	// below to trap any lease request and infer that it refers to the range we're
	// interested in.
	sc.TestingKnobs.DisableSplitQueue = true
	sc.TestingKnobs.LeaseRequestEvent = func(ts hlc.Timestamp) {
		val := leaseAcquisitionTrap.Load()
		if val == nil {
			return
		}
		trapCallback := val.(func(ts hlc.Timestamp))
		if trapCallback != nil {
			trapCallback(ts)
		}
	}
	mtc := &multiTestContext{storeConfig: &sc}
	defer mtc.Stop()
	mtc.Start(t, 1)

	// Send a read, to acquire a lease.
	getArgs := getArgs([]byte("a"))
	if _, err := client.SendWrapped(ctx, mtc.stores[0].TestSender(), getArgs); err != nil {
		t.Fatal(err)
	}

	preRepl1, err := mtc.stores[0].GetReplica(1)
	if err != nil {
		t.Fatal(err)
	}
	preRestartLease, _ := preRepl1.GetLease()

	mtc.manualClock.Increment(1E9)

	// Restart the mtc. Before we do that, we're installing a callback used to
	// assert that a new lease has been requested. The callback is installed
	// before the restart, as the lease might be requested at any time and for
	// many reasons by background processes, even before we send the read below.
	leaseAcquisitionCh := make(chan error)
	var once sync.Once
	leaseAcquisitionTrap.Store(func(_ hlc.Timestamp) {
		once.Do(func() {
			close(leaseAcquisitionCh)
		})
	})

	log.Info(ctx, "restarting")
	mtc.restart()

	// Send another read and check that the pre-existing lease has not been used.
	// Concretely, we check that a new lease is requested.
	if _, err := client.SendWrapped(ctx, mtc.stores[0].TestSender(), getArgs); err != nil {
		t.Fatal(err)
	}
	// Check that the Send above triggered a lease acquisition.
	select {
	case <-leaseAcquisitionCh:
	case <-time.After(time.Second):
		t.Fatalf("read did not acquire a new lease")
	}

	postRepl1, err := mtc.stores[0].GetReplica(1)
	if err != nil {
		t.Fatal(err)
	}
	postRestartLease, _ := postRepl1.GetLease()

	// Verify that not only is a new lease requested, it also gets a new sequence
	// number. This makes sure that previously proposed commands actually fail at
	// apply time.
	if preRestartLease.Sequence == postRestartLease.Sequence {
		t.Fatalf("lease was not replaced:\nprev: %v\nnow:  %v", preRestartLease, postRestartLease)
	}
}

// Test that a lease extension (a RequestLeaseRequest that doesn't change the
// lease holder) is not blocked by ongoing reads.
// The test relies on two things:
// 1) Lease extensions, unlike lease transfers, are not blocked by reads through
// their ReplicatedEvalResult.BlockReads.
// 2) Requests such as RequestLeaseRequest don't declare to touch the whole key
// span of the range, and thus don't conflict through the command queue with
// other reads.
func TestLeaseExtensionNotBlockedByRead(t *testing.T) {
	defer leaktest.AfterTest(t)()
	readBlocked := make(chan struct{})
	cmdFilter := func(fArgs storagebase.FilterArgs) *roachpb.Error {
		if fArgs.Hdr.UserPriority == 42 {
			// Signal that the read is blocked.
			readBlocked <- struct{}{}
			// Wait for read to be unblocked.
			<-readBlocked
		}
		return nil
	}
	srv, _, _ := serverutils.StartServer(t,
		base.TestServerArgs{
			Knobs: base.TestingKnobs{
				Store: &storage.StoreTestingKnobs{
					EvalKnobs: storagebase.BatchEvalTestingKnobs{
						TestingEvalFilter: cmdFilter,
					},
				},
			},
		})
	s := srv.(*server.TestServer)
	defer s.Stopper().Stop(context.TODO())

	store, err := s.GetStores().(*storage.Stores).GetStore(s.GetFirstStoreID())
	if err != nil {
		t.Fatal(err)
	}

	// Start a read and wait for it to block.
	key := roachpb.Key("a")
	errChan := make(chan error)
	go func() {
		getReq := roachpb.GetRequest{
			RequestHeader: roachpb.RequestHeader{
				Key: key,
			},
		}
		if _, pErr := client.SendWrappedWith(context.Background(), s.DB().NonTransactionalSender(),
			roachpb.Header{UserPriority: 42},
			&getReq); pErr != nil {
			errChan <- pErr.GoError()
		}
	}()

	select {
	case err := <-errChan:
		t.Fatal(err)
	case <-readBlocked:
		// Send the lease request.
		rKey, err := keys.Addr(key)
		if err != nil {
			t.Fatal(err)
		}
		repl := store.LookupReplica(rKey)
		if repl == nil {
			t.Fatalf("replica for key %s not found", rKey)
		}
		replDesc, found := repl.Desc().GetReplicaDescriptor(store.StoreID())
		if !found {
			t.Fatalf("replica descriptor for key %s not found", rKey)
		}

		curLease, _, err := s.GetRangeLease(context.TODO(), key)
		if err != nil {
			t.Fatal(err)
		}

		leaseReq := roachpb.RequestLeaseRequest{
			RequestHeader: roachpb.RequestHeader{
				Key: key,
			},
			Lease: roachpb.Lease{
				Start:      s.Clock().Now(),
				Expiration: s.Clock().Now().Add(time.Second.Nanoseconds(), 0).Clone(),
				Replica:    replDesc,
			},
			PrevLease: curLease,
		}
		if _, pErr := client.SendWrapped(context.Background(), s.DB().NonTransactionalSender(), &leaseReq); pErr != nil {
			t.Fatal(pErr)
		}
		// Unblock the read.
		readBlocked <- struct{}{}
	}
}

// LeaseInfo runs a LeaseInfoRequest using the specified server.
func LeaseInfo(
	t *testing.T,
	db *client.DB,
	rangeDesc roachpb.RangeDescriptor,
	readConsistency roachpb.ReadConsistencyType,
) roachpb.LeaseInfoResponse {
	leaseInfoReq := &roachpb.LeaseInfoRequest{
		RequestHeader: roachpb.RequestHeader{
			Key: rangeDesc.StartKey.AsRawKey(),
		},
	}
	reply, pErr := client.SendWrappedWith(context.Background(), db.NonTransactionalSender(), roachpb.Header{
		ReadConsistency: readConsistency,
	}, leaseInfoReq)
	if pErr != nil {
		t.Fatal(pErr)
	}
	return *(reply.(*roachpb.LeaseInfoResponse))
}

func TestLeaseInfoRequest(t *testing.T) {
	defer leaktest.AfterTest(t)()
	tc := testcluster.StartTestCluster(t, 3,
		base.TestClusterArgs{
			ReplicationMode: base.ReplicationManual,
		})
	defer tc.Stopper().Stop(context.TODO())

	kvDB0 := tc.Servers[0].DB()
	kvDB1 := tc.Servers[1].DB()

	key := []byte("a")
	rangeDesc, err := tc.LookupRange(key)
	if err != nil {
		t.Fatal(err)
	}

	rangeDesc, err = tc.AddReplicas(
		rangeDesc.StartKey.AsRawKey(), tc.Target(1), tc.Target(2),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(rangeDesc.Replicas) != 3 {
		t.Fatalf("expected 3 replicas, got %+v", rangeDesc.Replicas)
	}
	replicas := make([]roachpb.ReplicaDescriptor, 3)
	for i := 0; i < 3; i++ {
		var ok bool
		replicas[i], ok = rangeDesc.GetReplicaDescriptor(tc.Servers[i].GetFirstStoreID())
		if !ok {
			t.Fatalf("expected to find replica in server %d", i)
		}
	}

	// Lease should start on Server 0, since nobody told it to move.
	leaseHolderReplica := LeaseInfo(t, kvDB0, rangeDesc, roachpb.INCONSISTENT).Lease.Replica
	if leaseHolderReplica != replicas[0] {
		t.Fatalf("lease holder should be replica %+v, but is: %+v", replicas[0], leaseHolderReplica)
	}

	// Transfer the lease to Server 1 and check that LeaseInfoRequest gets the
	// right answer.
	err = tc.TransferRangeLease(rangeDesc, tc.Target(1))
	if err != nil {
		t.Fatal(err)
	}
	// An inconsistent LeaseInfoReqeust on the old lease holder should give us the
	// right answer immediately, since the old holder has definitely applied the
	// transfer before TransferRangeLease returned.
	leaseHolderReplica = LeaseInfo(t, kvDB0, rangeDesc, roachpb.INCONSISTENT).Lease.Replica
	if leaseHolderReplica != replicas[1] {
		t.Fatalf("lease holder should be replica %+v, but is: %+v",
			replicas[1], leaseHolderReplica)
	}

	// A read on the new lease holder does not necessarily succeed immediately,
	// since it might take a while for it to apply the transfer.
	testutils.SucceedsSoon(t, func() error {
		// We can't reliably do a CONSISTENT read here, even though we're reading
		// from the supposed lease holder, because this node might initially be
		// unaware of the new lease and so the request might bounce around for a
		// while (see #8816).
		leaseHolderReplica = LeaseInfo(t, kvDB1, rangeDesc, roachpb.INCONSISTENT).Lease.Replica
		if leaseHolderReplica != replicas[1] {
			return errors.Errorf("lease holder should be replica %+v, but is: %+v",
				replicas[1], leaseHolderReplica)
		}
		return nil
	})

	// Transfer the lease to Server 2 and check that LeaseInfoRequest gets the
	// right answer.
	err = tc.TransferRangeLease(rangeDesc, tc.Target(2))
	if err != nil {
		t.Fatal(err)
	}
	leaseHolderReplica = LeaseInfo(t, kvDB1, rangeDesc, roachpb.INCONSISTENT).Lease.Replica
	if leaseHolderReplica != replicas[2] {
		t.Fatalf("lease holder should be replica %+v, but is: %+v", replicas[2], leaseHolderReplica)
	}

	// TODO(andrei): test the side-effect of LeaseInfoRequest when there's no
	// active lease - the node getting the request is supposed to acquire the
	// lease. This requires a way to expire leases; the TestCluster probably needs
	// to use a mock clock.
}

// Test that an error encountered by a read-only "NonKV" command is not
// swallowed, and doesn't otherwise cause a panic.
// We had a bug cause by the fact that errors for these commands aren't passed
// through the epilogue returned by replica.beginCommands() and were getting
// swallowed.
func TestErrorHandlingForNonKVCommand(t *testing.T) {
	defer leaktest.AfterTest(t)()
	cmdFilter := func(fArgs storagebase.FilterArgs) *roachpb.Error {
		if fArgs.Hdr.UserPriority == 42 {
			return roachpb.NewErrorf("injected error")
		}
		return nil
	}
	srv, _, _ := serverutils.StartServer(t,
		base.TestServerArgs{
			Knobs: base.TestingKnobs{
				Store: &storage.StoreTestingKnobs{
					EvalKnobs: storagebase.BatchEvalTestingKnobs{
						TestingEvalFilter: cmdFilter,
					},
				},
			},
		})
	s := srv.(*server.TestServer)
	defer s.Stopper().Stop(context.TODO())

	// Send the lease request.
	key := roachpb.Key("a")
	leaseReq := roachpb.LeaseInfoRequest{
		RequestHeader: roachpb.RequestHeader{
			Key: key,
		},
	}
	_, pErr := client.SendWrappedWith(
		context.Background(),
		s.DB().NonTransactionalSender(),
		roachpb.Header{UserPriority: 42},
		&leaseReq,
	)
	if !testutils.IsPError(pErr, "injected error") {
		t.Fatalf("expected error %q, got: %s", "injected error", pErr)
	}
}

func TestRangeInfo(t *testing.T) {
	defer leaktest.AfterTest(t)()
	storeCfg := storage.TestStoreConfig(nil /* clock */)
	storeCfg.TestingKnobs.DisableMergeQueue = true
	mtc := &multiTestContext{
		storeConfig: &storeCfg,
	}
	defer mtc.Stop()
	mtc.Start(t, 2)

	// Up-replicate to two replicas.
	mtc.replicateRange(mtc.stores[0].LookupReplica(roachpb.RKeyMin).RangeID, 1)

	// Split the key space at key "a".
	splitKey := roachpb.RKey("a")
	splitArgs := adminSplitArgs(splitKey.AsRawKey())
	if _, pErr := client.SendWrapped(
		context.Background(), mtc.stores[0].TestSender(), splitArgs,
	); pErr != nil {
		t.Fatal(pErr)
	}

	// Get the replicas for each side of the split. This is done within
	// a SucceedsSoon loop to ensure the split completes.
	var lhsReplica0, lhsReplica1, rhsReplica0, rhsReplica1 *storage.Replica
	testutils.SucceedsSoon(t, func() error {
		lhsReplica0 = mtc.stores[0].LookupReplica(roachpb.RKeyMin)
		lhsReplica1 = mtc.stores[1].LookupReplica(roachpb.RKeyMin)
		rhsReplica0 = mtc.stores[0].LookupReplica(splitKey)
		rhsReplica1 = mtc.stores[1].LookupReplica(splitKey)
		if lhsReplica0 == rhsReplica0 || lhsReplica1 == rhsReplica1 {
			return errors.Errorf("replicas not post-split %v, %v, %v, %v",
				lhsReplica0, rhsReplica0, rhsReplica0, rhsReplica1)
		}
		return nil
	})
	lhsLease, _ := lhsReplica0.GetLease()
	rhsLease, _ := rhsReplica0.GetLease()

	// Verify range info is not set if unrequested.
	getArgs := getArgs(splitKey.AsRawKey())
	reply, pErr := client.SendWrapped(context.Background(), mtc.distSenders[0], getArgs)
	if pErr != nil {
		t.Fatal(pErr)
	}
	if len(reply.Header().RangeInfos) > 0 {
		t.Errorf("expected empty range infos if unrequested; got %v", reply.Header().RangeInfos)
	}

	// Verify range info on a get request.
	h := roachpb.Header{
		ReturnRangeInfo: true,
	}
	reply, pErr = client.SendWrappedWith(context.Background(), mtc.distSenders[0], h, getArgs)
	if pErr != nil {
		t.Fatal(pErr)
	}
	expRangeInfos := []roachpb.RangeInfo{
		{
			Desc:  *rhsReplica0.Desc(),
			Lease: rhsLease,
		},
	}
	if !reflect.DeepEqual(reply.Header().RangeInfos, expRangeInfos) {
		t.Errorf("on get reply, expected %+v; got %+v", expRangeInfos, reply.Header().RangeInfos)
	}

	// Verify range info on a put request.
	putArgs := putArgs(splitKey.AsRawKey(), []byte("foo"))
	reply, pErr = client.SendWrappedWith(context.Background(), mtc.distSenders[0], h, putArgs)
	if pErr != nil {
		t.Fatal(pErr)
	}
	if !reflect.DeepEqual(reply.Header().RangeInfos, expRangeInfos) {
		t.Errorf("on put reply, expected %+v; got %+v", expRangeInfos, reply.Header().RangeInfos)
	}

	// Verify range info on an admin request.
	adminArgs := &roachpb.AdminTransferLeaseRequest{
		RequestHeader: roachpb.RequestHeader{
			Key: splitKey.AsRawKey(),
		},
		Target: rhsLease.Replica.StoreID,
	}
	reply, pErr = client.SendWrappedWith(context.Background(), mtc.distSenders[0], h, adminArgs)
	if pErr != nil {
		t.Fatal(pErr)
	}
	if !reflect.DeepEqual(reply.Header().RangeInfos, expRangeInfos) {
		t.Errorf("on admin reply, expected %+v; got %+v", expRangeInfos, reply.Header().RangeInfos)
	}

	// Verify multiple range infos on a scan request.
	scanArgs := roachpb.ScanRequest{
		RequestHeader: roachpb.RequestHeader{
			Key:    keys.SystemMax,
			EndKey: roachpb.KeyMax,
		},
	}
	txn := roachpb.MakeTransaction("test", roachpb.KeyMin, 1, enginepb.SERIALIZABLE, mtc.clock.Now(), 0)
	h.Txn = &txn
	reply, pErr = client.SendWrappedWith(context.Background(), mtc.distSenders[0], h, &scanArgs)
	if pErr != nil {
		t.Fatal(pErr)
	}
	expRangeInfos = []roachpb.RangeInfo{
		{
			Desc:  *lhsReplica0.Desc(),
			Lease: lhsLease,
		},
		{
			Desc:  *rhsReplica0.Desc(),
			Lease: rhsLease,
		},
	}
	if !reflect.DeepEqual(reply.Header().RangeInfos, expRangeInfos) {
		t.Errorf("on scan reply, expected %+v; got %+v", expRangeInfos, reply.Header().RangeInfos)
	}

	// Verify multiple range infos and order on a reverse scan request.
	revScanArgs := roachpb.ReverseScanRequest{
		RequestHeader: roachpb.RequestHeader{
			Key:    keys.SystemMax,
			EndKey: roachpb.KeyMax,
		},
	}
	reply, pErr = client.SendWrappedWith(context.Background(), mtc.distSenders[0], h, &revScanArgs)
	if pErr != nil {
		t.Fatal(pErr)
	}
	expRangeInfos = []roachpb.RangeInfo{
		{
			Desc:  *rhsReplica0.Desc(),
			Lease: rhsLease,
		},
		{
			Desc:  *lhsReplica0.Desc(),
			Lease: lhsLease,
		},
	}
	if !reflect.DeepEqual(reply.Header().RangeInfos, expRangeInfos) {
		t.Errorf("on reverse scan reply, expected %+v; got %+v", expRangeInfos, reply.Header().RangeInfos)
	}

	// Change lease holders for both ranges and re-scan.
	for _, r := range []*storage.Replica{lhsReplica1, rhsReplica1} {
		replDesc, err := r.GetReplicaDescriptor()
		if err != nil {
			t.Fatal(err)
		}
		if err = mtc.dbs[0].AdminTransferLease(context.TODO(),
			r.Desc().StartKey.AsRawKey(), replDesc.StoreID); err != nil {
			t.Fatalf("unable to transfer lease to replica %s: %s", r, err)
		}
	}
	reply, pErr = client.SendWrappedWith(context.Background(), mtc.distSenders[0], h, &scanArgs)
	if pErr != nil {
		t.Fatal(pErr)
	}
	lhsLease, _ = lhsReplica1.GetLease()
	rhsLease, _ = rhsReplica1.GetLease()
	expRangeInfos = []roachpb.RangeInfo{
		{
			Desc:  *lhsReplica1.Desc(),
			Lease: lhsLease,
		},
		{
			Desc:  *rhsReplica1.Desc(),
			Lease: rhsLease,
		},
	}
	if !reflect.DeepEqual(reply.Header().RangeInfos, expRangeInfos) {
		t.Errorf("on scan reply, expected %+v; got %+v", expRangeInfos, reply.Header().RangeInfos)
	}
}

// TestDrainRangeRejection verifies that an attempt to transfer a range to a
// draining store fails.
func TestDrainRangeRejection(t *testing.T) {
	defer leaktest.AfterTest(t)()
	mtc := &multiTestContext{}
	defer mtc.Stop()
	mtc.Start(t, 2)

	repl, err := mtc.stores[0].GetReplica(1)
	if err != nil {
		t.Fatal(err)
	}

	drainingIdx := 1
	mtc.stores[drainingIdx].SetDraining(true)
	if err := repl.ChangeReplicas(
		context.Background(),
		roachpb.ADD_REPLICA,
		roachpb.ReplicationTarget{
			NodeID:  mtc.idents[drainingIdx].NodeID,
			StoreID: mtc.idents[drainingIdx].StoreID,
		},
		repl.Desc(),
		storagepb.ReasonRangeUnderReplicated,
		"",
	); !testutils.IsError(err, "store is draining") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSystemZoneConfigs(t *testing.T) {
	defer leaktest.AfterTest(t)()

	// This test is relatively slow and resource intensive. When run under
	// stressrace on a loaded machine (as in the nightly tests), sometimes the
	// SucceedsSoon conditions below take longer than the allotted time (#25273).
	if testing.Short() || testutils.NightlyStress() {
		t.Skip()
	}

	ctx := context.Background()
	tc := testcluster.StartTestCluster(t, 7, base.TestClusterArgs{
		ServerArgs: base.TestServerArgs{
			Knobs: base.TestingKnobs{
				Store: &storage.StoreTestingKnobs{
					// Disable LBS because when the scan is happening at the rate it's happening
					// below, it's possible that one of the system ranges trigger a split.
					DisableLoadBasedSplitting: true,
				},
			},
			// Scan like a bat out of hell to ensure replication and replica GC
			// happen in a timely manner.
			ScanInterval: 50 * time.Millisecond,
		},
	})
	defer tc.Stopper().Stop(ctx)
	log.Info(ctx, "TestSystemZoneConfig: test cluster started")

	expectedSystemRanges, err := tc.Servers[0].ExpectedInitialRangeCount()
	if err != nil {
		t.Fatal(err)
	}
	expectedUserRanges := tc.Servers[0].ExpectedInitialUserRangeCount()
	expectedSystemRanges -= expectedUserRanges
	systemNumReplicas := int(*config.DefaultSystemZoneConfig().NumReplicas)
	userNumReplicas := int(*config.DefaultZoneConfig().NumReplicas)
	expectedReplicas := expectedSystemRanges*systemNumReplicas + expectedUserRanges*userNumReplicas
	log.Infof(ctx, "TestSystemZoneConfig: expecting %d system ranges and %d user ranges",
		expectedSystemRanges, expectedUserRanges)
	log.Infof(ctx, "TestSystemZoneConfig: expected (%dx%d) + (%dx%d) = %d replicas total",
		expectedSystemRanges, systemNumReplicas, expectedUserRanges, userNumReplicas, expectedReplicas)

	waitForReplicas := func() error {
		var conflictingID roachpb.RangeID
		replicas := make(map[roachpb.RangeID]int)
		for _, s := range tc.Servers {
			if err := storage.IterateRangeDescriptors(ctx, s.Engines()[0], func(desc roachpb.RangeDescriptor) (bool, error) {
				if existing, ok := replicas[desc.RangeID]; ok && existing != len(desc.Replicas) {
					conflictingID = desc.RangeID
				}
				replicas[desc.RangeID] = len(desc.Replicas)
				return false, nil
			}); err != nil {
				return err
			}
		}
		if conflictingID != 0 {
			return fmt.Errorf("not all replicas agree on the range descriptor for r%d", conflictingID)
		}
		var totalReplicas int
		for _, count := range replicas {
			totalReplicas += count
		}
		if totalReplicas != expectedReplicas {
			return fmt.Errorf("got %d replicas, want %d; details: %+v", totalReplicas, expectedReplicas, replicas)
		}
		return nil
	}

	// Wait until we're down to the expected number of replicas. This is
	// effectively waiting on replica GC to kick in to destroy any replicas that
	// got removed during rebalancing of the initial ranges, since the testcluster
	// waits until nothing is underreplicated but not until all rebalancing has
	// settled down.
	testutils.SucceedsSoon(t, waitForReplicas)
	log.Info(ctx, "TestSystemZoneConfig: initial replication succeeded")

	// Update the meta zone config to have more replicas and expect the number
	// of replicas to go up accordingly after running all replicas through the
	// replicate queue.
	sqlDB := sqlutils.MakeSQLRunner(tc.ServerConn(0))
	sqlutils.SetZoneConfig(t, sqlDB, "RANGE meta", "num_replicas: 7")
	expectedReplicas += 2
	testutils.SucceedsSoon(t, waitForReplicas)
	log.Info(ctx, "TestSystemZoneConfig: up-replication of meta ranges succeeded")

	// Do the same thing, but down-replicating the timeseries range.
	sqlutils.SetZoneConfig(t, sqlDB, "RANGE timeseries", "num_replicas: 1")
	expectedReplicas -= 2
	testutils.SucceedsSoon(t, waitForReplicas)
	log.Info(ctx, "TestSystemZoneConfig: down-replication of timeseries ranges succeeded")

	// Finally, verify the system ranges. Note that in a new cluster there are
	// two system ranges, which we have to take into account here.
	sqlutils.SetZoneConfig(t, sqlDB, "RANGE system", "num_replicas: 7")
	expectedReplicas += 6
	testutils.SucceedsSoon(t, waitForReplicas)
	log.Info(ctx, "TestSystemZoneConfig: up-replication of system ranges succeeded")
}

func TestClearRange(t *testing.T) {
	defer leaktest.AfterTest(t)()

	ctx := context.Background()
	stopper := stop.NewStopper()
	defer stopper.Stop(ctx)
	store := createTestStoreWithConfig(t, stopper, storage.TestStoreConfig(nil))

	clearRange := func(start, end roachpb.Key) {
		t.Helper()
		if _, err := client.SendWrapped(ctx, store.DB().NonTransactionalSender(), &roachpb.ClearRangeRequest{
			RequestHeader: roachpb.RequestHeader{
				Key:    start,
				EndKey: end,
			},
		}); err != nil {
			t.Fatal(err)
		}
	}

	verifyKeysWithPrefix := func(prefix roachpb.Key, expectedKeys []roachpb.Key) {
		t.Helper()
		start := engine.MakeMVCCMetadataKey(prefix)
		end := engine.MakeMVCCMetadataKey(prefix.PrefixEnd())
		kvs, err := engine.Scan(store.Engine(), start, end, 0 /* maxRows */)
		if err != nil {
			t.Fatal(err)
		}
		var actualKeys []roachpb.Key
		for _, kv := range kvs {
			actualKeys = append(actualKeys, kv.Key.Key)
		}
		if !reflect.DeepEqual(expectedKeys, actualKeys) {
			t.Fatalf("expected %v, but got %v", expectedKeys, actualKeys)
		}
	}

	rng, _ := randutil.NewPseudoRand()

	// Write four keys with values small enough to use individual deletions
	// (sm1-sm4) and four keys with values large enough to require a range
	// deletion tombstone (lg1-lg4).
	sm, sm1, sm2, sm3 := roachpb.Key("sm"), roachpb.Key("sm1"), roachpb.Key("sm2"), roachpb.Key("sm3")
	lg, lg1, lg2, lg3 := roachpb.Key("lg"), roachpb.Key("lg1"), roachpb.Key("lg2"), roachpb.Key("lg3")
	for _, key := range []roachpb.Key{sm1, sm2, sm3} {
		if err := store.DB().Put(ctx, key, "sm-val"); err != nil {
			t.Fatal(err)
		}
	}
	for _, key := range []roachpb.Key{lg1, lg2, lg3} {
		if err := store.DB().Put(
			ctx, key, randutil.RandBytes(rng, batcheval.ClearRangeBytesThreshold),
		); err != nil {
			t.Fatal(err)
		}
	}
	verifyKeysWithPrefix(sm, []roachpb.Key{sm1, sm2, sm3})
	verifyKeysWithPrefix(lg, []roachpb.Key{lg1, lg2, lg3})

	// Verify that a ClearRange request from [sm1, sm3) removes sm1 and sm2.
	clearRange(sm1, sm3)
	verifyKeysWithPrefix(sm, []roachpb.Key{sm3})

	// Verify that a ClearRange request from [lg1, lg3) removes lg1 and lg2.
	clearRange(lg1, lg3)
	verifyKeysWithPrefix(lg, []roachpb.Key{lg3})

	// Verify that only the large ClearRange request used a range deletion
	// tombstone by checking for the presence of a suggested compaction.
	verifyKeysWithPrefix(keys.LocalStoreSuggestedCompactionsMin,
		[]roachpb.Key{keys.StoreSuggestedCompactionKey(lg1, lg3)})
}
