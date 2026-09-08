package ip

import "testing"

func TestBuildHarnessDefaultsToTenMillionYearlyTarget(t *testing.T) {
	harness, err := BuildHarness(Options{
		TargetYearlyRevenue: 10000000,
		Months:              12,
		AverageOrderValue:   20000,
		LeadToCustomerRate:  0.02,
	})
	if err != nil {
		t.Fatal(err)
	}

	if harness.TargetYearlyRevenue != 10000000 {
		t.Fatalf("expected 10000000 yearly target, got %d", harness.TargetYearlyRevenue)
	}
	if harness.FunnelTargets.MonthlyRevenueTarget != 833334 {
		t.Fatalf("unexpected monthly revenue target: %d", harness.FunnelTargets.MonthlyRevenueTarget)
	}
	if harness.FunnelTargets.MonthlyCustomers != 42 {
		t.Fatalf("unexpected monthly customers: %d", harness.FunnelTargets.MonthlyCustomers)
	}
	if harness.FunnelTargets.MonthlyLeads != 2100 {
		t.Fatalf("unexpected monthly leads: %d", harness.FunnelTargets.MonthlyLeads)
	}

	var total int64
	for _, stream := range harness.RevenueStreams {
		total += stream.YearlyTarget
	}
	if total != harness.TargetYearlyRevenue {
		t.Fatalf("expected stream total to match target, got %d", total)
	}
}

func TestBuildHarnessRejectsInvalidConversionRate(t *testing.T) {
	_, err := BuildHarness(Options{
		TargetYearlyRevenue: 10000000,
		Months:              12,
		AverageOrderValue:   20000,
		LeadToCustomerRate:  0,
	})
	if err == nil {
		t.Fatal("expected invalid conversion rate error")
	}
}
