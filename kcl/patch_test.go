package kcl

// Fork patch, not upstream: the daemon sends these fields, the stock library
// drops them. These tests fail on the upstream code and guard the patch.

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// The audit trail keys S3 objects on subSequenceNumber and stores
// approximateArrivalTimestamp. The stock Record drops both, so this fails on
// the upstream library: it guards the patch against an update that loses it.
func TestRecordKeepsSubSequenceAndArrival(t *testing.T) {
	// A record as the MultiLangDaemon (KCL 3) sends it.
	raw := `{"data":"dGVzdA==","partitionKey":"pk","sequenceNumber":"4967",` +
		`"subSequenceNumber":3,"approximateArrivalTimestamp":1791000100123}`
	var r Record
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatal(err)
	}
	if r.SubSequenceNumber != 3 || r.ApproximateArrivalTimestamp != 1791000100123 {
		t.Fatalf("got subSequenceNumber %d, approximateArrivalTimestamp %d",
			r.SubSequenceNumber, r.ApproximateArrivalTimestamp)
	}
	if string(r.Data) != "test" {
		t.Fatalf("data %q", r.Data)
	}
}

// The daemon sends the lease's checkpoint with initialize; it must reach the
// processor, so a restart can be reconciled with what was already persisted.
func TestRun_InitializeCarriesCheckpoint(t *testing.T) {
	for _, tc := range []struct {
		line string
		seq  string
		sub  int64
	}{
		{`{"action": "initialize", "shardId": "s", "sequenceNumber": "49612345678901234567890", "subSequenceNumber": 3}`, "49612345678901234567890", 3},
		{`{"action": "initialize", "shardId": "s", "sequenceNumber": "TRIM_HORIZON", "subSequenceNumber": null}`, "TRIM_HORIZON", 0},
		{`{"action": "initialize", "shardId": "s"}`, "", 0},
	} {
		mProcessor := &mockProcessor{}
		k := &kclProcess{
			recordProcessor: mProcessor,
			logger:          defaultLogger,
			reader:          bufio.NewReader(strings.NewReader(tc.line + "\n")),
			writer:          bufio.NewWriter(&bytes.Buffer{}),
		}
		_ = k.Run() // ends at EOF
		got := mProcessor.initializeCall
		if got == nil || got.SequenceNumber != tc.seq || got.SubSequenceNumber != tc.sub {
			t.Errorf("%s: initialize input = %+v, want %s/%d", tc.line, got, tc.seq, tc.sub)
		}
	}
}
