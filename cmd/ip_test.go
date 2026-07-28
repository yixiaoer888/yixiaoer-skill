package cmd

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestIPHarnessCommandWritesStableJSONEnvelope(t *testing.T) {
	cmd := newIPHarnessCmd()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	t.Cleanup(resetIPHarnessFlags)

	if err := cmd.Flags().Set("target-yearly-revenue", "10000000"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("months", "12"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatal(err)
	}

	var envelope struct {
		OK     bool   `json:"ok"`
		Action string `json:"action"`
		Data   struct {
			Project             string `json:"project"`
			TargetYearlyRevenue int64  `json:"targetYearlyRevenue"`
			FunnelTargets       struct {
				MonthlyRevenueTarget int64 `json:"monthlyRevenueTarget"`
			} `json:"funnelTargets"`
		} `json:"data"`
	}
	if err := json.Unmarshal(out.Bytes(), &envelope); err != nil {
		t.Fatal(err)
	}
	if !envelope.OK {
		t.Fatal("expected ok envelope")
	}
	if envelope.Action != "ip.harness" {
		t.Fatalf("unexpected action: %s", envelope.Action)
	}
	if envelope.Data.Project == "" {
		t.Fatal("expected project name")
	}
	if envelope.Data.TargetYearlyRevenue != 10000000 {
		t.Fatalf("unexpected yearly target: %d", envelope.Data.TargetYearlyRevenue)
	}
	if envelope.Data.FunnelTargets.MonthlyRevenueTarget != 833334 {
		t.Fatalf("unexpected monthly target: %d", envelope.Data.FunnelTargets.MonthlyRevenueTarget)
	}
}

func resetIPHarnessFlags() {
	ipProjectName = ""
	ipTargetYearlyRevenue = 0
	ipMonths = 0
	ipAverageOrderValue = 0
	ipLeadToCustomerRate = 0
}
