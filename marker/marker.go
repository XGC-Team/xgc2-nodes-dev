package marker

import (
	"context"
	"errors"
	"strings"

	"github.com/lxk36/xgc2-orchestration-core/kernel/canonicaljson"
	protocol "github.com/lxk36/xgc2-orchestration-core/kernel/node"
	"github.com/lxk36/xgc2-orchestration-core/sdk/go/contracts"
)

const packageDigest = "sha256:05dfac5608ec50518fda75dd93ab31ba97694cac56fd3398a5710a37c5cbed29"

type Executor struct{ descriptor contracts.NodeDescriptor }

func New() *Executor {
	descriptor := contracts.NodeDescriptor{
		SchemaVersion: protocol.DescriptorSchemaVersion, TypeRef: "xgc.dev.marker-normalize/v1", DisplayName: "Normalize product marker",
		PackageRef: "xgc2-nodes-dev", PackageDigest: packageDigest,
		InputSchema: contracts.Schema{Type: contracts.TypeObject, Properties: map[string]contracts.Schema{
			"product": {Type: contracts.TypeString}, "location": {Type: contracts.TypeString}, "description": {Type: contracts.TypeString}, "severity": {Type: contracts.TypeString},
		}, Required: []string{"product", "location", "description", "severity"}},
		OutputSchema: contracts.Schema{Type: contracts.TypeObject, Properties: map[string]contracts.Schema{
			"markerId": {Type: contracts.TypeString}, "product": {Type: contracts.TypeString}, "location": {Type: contracts.TypeString}, "description": {Type: contracts.TypeString}, "severity": {Type: contracts.TypeString},
		}, Required: []string{"markerId", "product", "location", "description", "severity"}},
		Mode: contracts.NodePure, Determinism: contracts.NodeDeterministic, MaxInputBytes: 65536, MaxOutputBytes: 65536,
	}
	descriptor.DescriptorDigest, _ = protocol.DescriptorDigest(descriptor)
	return &Executor{descriptor: descriptor}
}

func (executor *Executor) Descriptor() contracts.NodeDescriptor { return executor.descriptor }

func (executor *Executor) Execute(_ context.Context, request contracts.NodeInvocationRequest) (contracts.NodeResult, error) {
	product, pOK := request.Input["product"].(string)
	location, lOK := request.Input["location"].(string)
	description, dOK := request.Input["description"].(string)
	severity, sOK := request.Input["severity"].(string)
	if !pOK || !lOK || !dOK || !sOK || strings.TrimSpace(product) == "" || strings.TrimSpace(location) == "" || strings.TrimSpace(description) == "" {
		return contracts.NodeResult{}, errors.New("marker product, location, or description is invalid")
	}
	severity = strings.ToLower(strings.TrimSpace(severity))
	if severity != "low" && severity != "medium" && severity != "high" && severity != "critical" {
		return contracts.NodeResult{}, errors.New("marker severity is invalid")
	}
	normalized := map[string]any{"product": strings.TrimSpace(product), "location": strings.TrimSpace(location), "description": strings.TrimSpace(description), "severity": severity}
	identity, err := canonicaljson.DigestValue(normalized)
	if err != nil {
		return contracts.NodeResult{}, err
	}
	normalized["markerId"] = "marker-" + identity[len("sha256:"):]
	digest, err := canonicaljson.DigestValue(normalized)
	if err != nil {
		return contracts.NodeResult{}, err
	}
	return contracts.NodeResult{Status: contracts.NodeResultSucceeded, Output: normalized, OutputDigest: digest, EvidenceDigest: digest}, nil
}
