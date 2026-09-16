package main

import (
	"fmt"
	"os"

	"github.com/DeBoX85/Cloud-Assess/internal/branding"
	"github.com/spf13/cobra"
)

func main() {
	brand := branding.Default()

	root := &cobra.Command{
		Use:   brand.CLIName,
		Short: brand.ReportTitle,
		Long:  fmt.Sprintf("%s is an Azure cloud assessment toolkit.", brand.ProductName),
	}

	root.AddCommand(&cobra.Command{
		Use:   "scan",
		Short: "Assess Azure resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("scan engine is not implemented yet")
		},
	})

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
