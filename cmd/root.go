package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
	"github.com/spf13/cobra"
)

var (
	extensions  []string
	projectRoot string
	statistics  bool
)
var useEffectCount int
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
				analyzeFile(path)
			}
			return nil
		})
		if err != nil {
			fmt.Println("An error has occured", err)
		}
		fmt.Println("total useEffects found", useEffectCount)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringSliceVar(&extensions, "files", []string{".tsx", ".jsx"}, "File extensions to scan")
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

func analyzeFile(path string) {
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to read %s: %v\n", path, err)
		return
	}

	parser := sitter.NewParser()
	parser.SetLanguage(typescript.GetLanguage())
	tree, _ := parser.ParseCtx(context.Background(), nil, source)
	root := tree.RootNode()

	findUseEffectCalls(root, source, path)
}

func findUseEffectCalls(node *sitter.Node, source []byte, path string) {
	if node == nil {
		return
	}

	if node.Type() == "call_expression" {
		funcNode := node.ChildByFieldName("function")
		if funcNode != nil && funcNode.Type() == "identifier" {
			funcName := funcNode.Content(source)
			if funcName == "useEffect" {
				useEffectCount++
				argsNode := node.ChildByFieldName("arguments")
				if argsNode != nil && argsNode.ChildCount() > 0 {
					if argsNode.ChildCount() > 1 {
						startByte := node.StartByte()
						endByte := node.EndByte()
						useEffectCall := source[startByte:endByte]

						fmt.Println("UseEffect", string(useEffectCall))
					}
				}
				fmt.Println("Found at", path)
				fmt.Println("--------------------------------------------------------------------")
			}
		}
	}

	for i := range int(node.ChildCount()) {
		findUseEffectCalls(node.Child(i), source, path)
	}
}

// func analyzeUseEffect(point sitter.Point)

// useEffect(() = > {

// }, [])
