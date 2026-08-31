package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yixiaoer/yixiaoer-skill/internal/output"
	ipflow "github.com/yixiaoer/yixiaoer-skill/internal/workflows/ip"
)

var (
	ipProjectName         string
	ipTargetYearlyRevenue int64
	ipMonths              int
	ipAverageOrderValue   int64
	ipLeadToCustomerRate  float64
)

func init() {
	rootCmd.AddCommand(newIPCmd())
}

func newIPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ip",
		Short: "自媒体 IP 项目规划工具",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newIPHarnessCmd())
	return cmd
}

func newIPHarnessCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "harness",
		Short: "生成年度营收目标的自媒体 IP 运营骨架",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runIPHarness(cmd)
		},
	}
	cmd.Flags().StringVar(&ipProjectName, "project", "media-ip-10m-arr", "project name")
	cmd.Flags().Int64Var(&ipTargetYearlyRevenue, "target-yearly-revenue", 10000000, "target yearly revenue")
	cmd.Flags().IntVar(&ipMonths, "months", 12, "planning months")
	cmd.Flags().Int64Var(&ipAverageOrderValue, "average-order-value", 20000, "average order value")
	cmd.Flags().Float64Var(&ipLeadToCustomerRate, "lead-to-customer-rate", 0.02, "lead to customer conversion rate")
	return cmd
}

func runIPHarness(cmd *cobra.Command) error {
	harness, err := ipflow.BuildHarness(ipflow.Options{
		Project:             ipProjectName,
		TargetYearlyRevenue: ipTargetYearlyRevenue,
		Months:              ipMonths,
		AverageOrderValue:   ipAverageOrderValue,
		LeadToCustomerRate:  ipLeadToCustomerRate,
	})
	if err != nil {
		return err
	}
	return output.Success(cmd.OutOrStdout(), "ip.harness", harness)
}
