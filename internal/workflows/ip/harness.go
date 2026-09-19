package ip

import (
	"math"

	"github.com/yixiaoer/yixiaoer-skill/internal/yxerrors"
)

type RevenueStream struct {
	Name          string `json:"name"`
	SharePercent  int    `json:"sharePercent"`
	YearlyTarget  int64  `json:"yearlyTarget"`
	MonthlyTarget int64  `json:"monthlyTarget"`
}

type ContentCadence struct {
	ShortVideosPerWeek  int `json:"shortVideosPerWeek"`
	LongFormPerWeek     int `json:"longFormPerWeek"`
	LiveSessionsPerWeek int `json:"liveSessionsPerWeek"`
	LeadMagnetsPerMonth int `json:"leadMagnetsPerMonth"`
}

type FunnelTargets struct {
	MonthlyRevenueTarget int64 `json:"monthlyRevenueTarget"`
	MonthlyLeads         int64 `json:"monthlyLeads"`
	MonthlySalesCalls    int64 `json:"monthlySalesCalls"`
	MonthlyCustomers     int64 `json:"monthlyCustomers"`
	AverageOrderValue    int64 `json:"averageOrderValue"`
}

type Harness struct {
	Project             string          `json:"project"`
	TargetYearlyRevenue int64           `json:"targetYearlyRevenue"`
	Months              int             `json:"months"`
	RevenueStreams      []RevenueStream `json:"revenueStreams"`
	ContentCadence      ContentCadence  `json:"contentCadence"`
	FunnelTargets       FunnelTargets   `json:"funnelTargets"`
	WeeklyOperatingLoop []string        `json:"weeklyOperatingLoop"`
	SuccessCriteria     []string        `json:"successCriteria"`
}

type Options struct {
	Project             string
	TargetYearlyRevenue int64
	Months              int
	AverageOrderValue   int64
	LeadToCustomerRate  float64
}

func BuildHarness(options Options) (Harness, error) {
	if options.Project == "" {
		options.Project = "media-ip-10m-arr"
	}
	if options.TargetYearlyRevenue <= 0 {
		return Harness{}, yxerrors.Usage("target yearly revenue must be greater than zero", options.TargetYearlyRevenue).
			WithHint("请使用 --target-yearly-revenue 传入大于 0 的年度营收目标。")
	}
	if options.Months <= 0 {
		return Harness{}, yxerrors.Usage("months must be greater than zero", options.Months).
			WithHint("请使用 --months 传入大于 0 的周期月数。")
	}
	if options.AverageOrderValue <= 0 {
		return Harness{}, yxerrors.Usage("average order value must be greater than zero", options.AverageOrderValue).
			WithHint("请使用 --average-order-value 传入大于 0 的客单价。")
	}
	if options.LeadToCustomerRate <= 0 || options.LeadToCustomerRate > 1 {
		return Harness{}, yxerrors.Usage("lead to customer rate must be in (0, 1]", options.LeadToCustomerRate).
			WithHint("请使用 --lead-to-customer-rate 传入 0 到 1 之间的转化率。")
	}

	monthlyRevenueTarget := ceilDiv(options.TargetYearlyRevenue, int64(options.Months))
	monthlyCustomers := ceilDiv(monthlyRevenueTarget, options.AverageOrderValue)
	monthlyLeads := int64(math.Ceil(float64(monthlyCustomers) / options.LeadToCustomerRate))

	return Harness{
		Project:             options.Project,
		TargetYearlyRevenue: options.TargetYearlyRevenue,
		Months:              options.Months,
		RevenueStreams: []RevenueStream{
			stream("sponsorship", 40, options.TargetYearlyRevenue, options.Months),
			stream("course", 30, options.TargetYearlyRevenue, options.Months),
			stream("consulting", 20, options.TargetYearlyRevenue, options.Months),
			stream("community", 10, options.TargetYearlyRevenue, options.Months),
		},
		ContentCadence: ContentCadence{
			ShortVideosPerWeek:  14,
			LongFormPerWeek:     2,
			LiveSessionsPerWeek: 1,
			LeadMagnetsPerMonth: 2,
		},
		FunnelTargets: FunnelTargets{
			MonthlyRevenueTarget: monthlyRevenueTarget,
			MonthlyLeads:         monthlyLeads,
			MonthlySalesCalls:    ceilDiv(monthlyLeads, 4),
			MonthlyCustomers:     monthlyCustomers,
			AverageOrderValue:    options.AverageOrderValue,
		},
		WeeklyOperatingLoop: []string{
			"research audience pain and monetizable demand",
			"publish acquisition content across selected channels",
			"route engaged audience into lead magnets",
			"review funnel metrics and update next week's content bets",
		},
		SuccessCriteria: []string{
			"revenue stream targets sum to yearly target",
			"monthly funnel target is derived from revenue target and average order value",
			"content cadence creates at least one weekly conversion moment",
		},
	}, nil
}

func stream(name string, sharePercent int, targetYearlyRevenue int64, months int) RevenueStream {
	yearlyTarget := targetYearlyRevenue * int64(sharePercent) / 100
	return RevenueStream{
		Name:          name,
		SharePercent:  sharePercent,
		YearlyTarget:  yearlyTarget,
		MonthlyTarget: ceilDiv(yearlyTarget, int64(months)),
	}
}

func ceilDiv(value int64, divisor int64) int64 {
	return (value + divisor - 1) / divisor
}
