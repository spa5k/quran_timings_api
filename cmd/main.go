package main

import (
	"fmt"
	"log"

	"github.com/spa5k/quran_timings/cmd/ayah"
	"github.com/spa5k/quran_timings/cmd/timings"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "quran_data"}
	rootCmd.AddCommand(
		makeCmd("all", "Run all data import functions sequentially", runAll),
		makeCmd("timings", "Fetch and insert timings data", timings.FetchQuranTimingReciters),
		makeCmd("ayah", "Fetch and insert ayah timings data", ayah.AyahTimingsPerReciter),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func makeCmd(name, description string, action func()) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: description,
		Run: func(cmd *cobra.Command, args []string) {
			action()
		},
	}
}

func runAll() {
	fmt.Println("Running all data import functions sequentially...")
	timings.FetchQuranTimingReciters()
	fmt.Println("All data import functions completed.")
}
