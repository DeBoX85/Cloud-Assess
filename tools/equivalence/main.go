package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/DeBoX85/Cloud-Assess/internal/equivalence"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	flags := flag.NewFlagSet("equivalence", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	var referencePath string
	var targetPath string
	var outputPath string
	flags.StringVar(&referencePath, "reference", "", "Pinned-reference JSON report")
	flags.StringVar(&targetPath, "target", "", "Cloud Assess canonical JSON report")
	flags.StringVar(&outputPath, "output", "", "Optional path for the JSON diff report")

	if err := flags.Parse(args); err != nil {
		return 2
	}
	if referencePath == "" || targetPath == "" {
		fmt.Fprintln(os.Stderr, "both --reference and --target are required")
		return 2
	}

	referenceFile, err := os.Open(referencePath) //nolint:gosec // caller explicitly supplies comparison input
	if err != nil {
		fmt.Fprintf(os.Stderr, "open reference report: %v\n", err)
		return 2
	}
	defer referenceFile.Close()

	targetFile, err := os.Open(targetPath) //nolint:gosec // caller explicitly supplies comparison input
	if err != nil {
		fmt.Fprintf(os.Stderr, "open target report: %v\n", err)
		return 2
	}
	defer targetFile.Close()

	reference, err := equivalence.LoadReference(referenceFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	target, err := equivalence.LoadTarget(targetFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	report := equivalence.Compare(reference, target)
	encoded, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "encode equivalence report: %v\n", err)
		return 2
	}
	encoded = append(encoded, '\n')

	if outputPath != "" {
		if err := os.WriteFile(outputPath, encoded, 0o600); err != nil {
			fmt.Fprintf(os.Stderr, "write equivalence report: %v\n", err)
			return 2
		}
	} else {
		if _, err := os.Stdout.Write(encoded); err != nil {
			fmt.Fprintf(os.Stderr, "write equivalence report: %v\n", err)
			return 2
		}
	}

	if !report.Equivalent {
		return 1
	}
	return 0
}
