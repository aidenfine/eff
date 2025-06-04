package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aidenfine/eff/src/analyzer"
	"github.com/spf13/cobra"
)

var (
	extensions  []string
	projectRoot string
	statistics  bool
)
var rootCmd = &cobra.Command{
	Use:   "eff",
	Short: "Analyze React useEffect usage",
	Run: func(cmd *cobra.Command, args []string) {
		if projectRoot == "" {
			projectRoot = "."
		}

		err := filepath.WalkDir(projectRoot, func(path string, d os.DirEntry, err error) error {
			if isIgnoredPath(path) {
			} else if isMatchingExtension(path) {

				if analyzer.AnalyzeFile(path) {
					fmt.Println(path, " is a barrel export")
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println("An error has occured", err)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringSliceVar(&extensions, "files", []string{".tsx", ".jsx", ".ts", ".js"}, "File extensions to scan")
	rootCmd.Flags().StringVar(&projectRoot, "dir", ".", "Root directory to scan")
	rootCmd.Flags().BoolVar(&statistics, "statistics", false, "Show statistics")
}

func isIgnoredPath(path string) bool {
	return strings.Contains(path, "node_modules")
}

func isMatchingExtension(path string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}
