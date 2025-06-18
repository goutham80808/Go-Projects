package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/goutham80808/valiform/internal/config"
	"github.com/goutham80808/valiform/internal/reader"
	"github.com/goutham80808/valiform/internal/validator"
	"github.com/goutham80808/valiform/internal/writer"
)

// We'll store the values from the command-line flags in these variables.
var (
	configFile string
	inputFile  string
	outputFile string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "valiform -c <config_file> -i <input_file>",
		Short: "Valiform is a CLI tool to validate and reformat data files.",
		Long: `Valiform is a flexible command-line tool that validates structured data files
(like CSV or JSON) against a user-defined set of rules in a YAML file.`,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("--- Valiform Starting ---")
			fmt.Printf("Config file: %s\n", configFile)
			fmt.Printf("Input file: %s\n", inputFile)
			fmt.Printf("Output file: %s\n", outputFile)
			fmt.Println("-------------------------")

			// 1. Load configuration
			rules, err := config.Load(configFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully loaded rules for %s files.\n", rules.FileType)

			// 2. Read input data
			// For now, we only support CSV, but this is where you'd add a switch for other types.
			var records []validator.Record
			if rules.FileType == "csv" {
				records, err = reader.ReadCSV(inputFile, rules.HasHeader)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
					os.Exit(1)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Unsupported file type '%s' in config\n", rules.FileType)
				os.Exit(1)
			}
			fmt.Printf("Read %d records from input file.\n", len(records))

			// 3. Validate records
			var validRecords []validator.Record
			var validationErrors []validator.ValidationError

			for i, record := range records {
				rowNum := i + 2
				if !rules.HasHeader {
					rowNum = i + 1
				}

				errors := validator.ValidateRecord(record, rules.Fields, rowNum)
				if len(errors) > 0 {
					validationErrors = append(validationErrors, errors...)
				} else {
					validRecords = append(validRecords, record)
				}
			}
			// 4. Write valid data and report results
			fmt.Println("--- Validation Complete ---")

			// Write the valid records to the output file if any exist.
			if len(validRecords) > 0 {
				err := writer.WriteJSON(outputFile, validRecords)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("Successfully wrote %d valid records to %s\n", len(validRecords), outputFile)
			} else {
				fmt.Println("No valid records found to write.")
			}
			if len(validationErrors) > 0 {
				fmt.Fprintf(os.Stderr, "\nFound %d validation errors:\n", len(validationErrors))
				for _, e := range validationErrors {
					fmt.Fprintln(os.Stderr, "- "+e.Error())
				}
				os.Exit(1)
			}

			fmt.Println("\nAll records processed successfully.")
		},
	}

	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the validation rules YAML file (required)")
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Path to the input data file (required)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "valid_output.json", "Path for the formatted output file")

	rootCmd.MarkFlagRequired("config")
	rootCmd.MarkFlagRequired("input")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
