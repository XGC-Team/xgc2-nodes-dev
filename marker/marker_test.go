package marker

import (
	"context"
	"testing"
	"time"

	"github.com/lxk36/xgc2-orchestration-core/conformance/nodepack"
	"github.com/lxk36/xgc2-orchestration-core/kernel/canonicaljson"
	"github.com/lxk36/xgc2-orchestration-core/sdk/go/contracts"
	nodesdk "github.com/lxk36/xgc2-orchestration-core/sdk/go/node"
)

func TestNodePackConformance(t *testing.T) {
	executor := New()
	input := map[string]any{"product": "XGC2", "location": "experiment panel", "description": "  stop button needs evidence  ", "severity": "HIGH"}
	digest, _ := canonicaljson.DigestValue(input)
	t0 := time.Date(2026, 8, 13, 2, 0, 0, 0, time.UTC)
	request := contracts.NodeInvocationRequest{InvocationID: "inv-1", RunID: "run-1", NodeID: "marker", TypeRef: executor.Descriptor().TypeRef, DescriptorDigest: executor.Descriptor().DescriptorDigest, AttemptID: "att-1", AttemptOrdinal: 1, Input: input, InputDigest: digest, RequestedAt: t0, Deadline: t0.Add(time.Minute)}
	report, err := nodepack.Validate(context.Background(), nodepack.Suite{PackageRef: "xgc2-nodes-dev", Executors: []nodesdk.Executor{executor}, Cases: []nodepack.Case{{Name: "normalize panel marker", Executor: executor, Request: request, ExpectedStatus: contracts.NodeResultSucceeded}}})
	if err != nil || report.DescriptorCount != 1 {
		t.Fatalf("report = %#v, err %v", report, err)
	}
}
