package analyzer

// what we could also do here is keep track of a score and weight
// for example if file name is index.ts, index.tsx, ...  higher chance of barrel import file increase weight

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aidenfine/eff/src/models"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

func AnalyzeFile(path string) models.FileFeedback {
	barrelExportReport := models.BarrelExportReport{
		Score: 0,
	}
	deadCodeReport := models.DeadCodeReport{
		Score: 0,
	}

	source, err := os.ReadFile(path)
	if err != nil {
		log.Panicln("Cannot read file", path)

	}
	parser := sitter.NewParser()
	parser.SetLanguage(typescript.GetLanguage())
	tree, _ := parser.ParseCtx(context.Background(), nil, source)
	root := tree.RootNode()
	fmt.Println(root, "root")
	findUseEffectCalls(root, source)
	if isBarrelFile(root, source) {
		fmt.Println("IS BARREL EXPORT")
		barrelExportReport.Score = 100
	}

	mainReport := models.FileFeedback{
		Score:               0,
		Path:                path,
		BarrelExportResults: barrelExportReport,
		DeadCodeResults:     deadCodeReport,
	}
	return mainReport

}
func isBarrelFile(root *sitter.Node, source []byte) bool {
	exportStatements := 0
	totalStatements := 0

	for i := range int(root.NamedChildCount()) {
		child := root.NamedChild(i)
		nodeType := child.Type()

		if nodeType == "export_all_statement" || nodeType == "export_statement" {
			content := child.Content(source)
			if strings.Contains(content, "from") {
				exportStatements++
			}
		}

		if nodeType != "comment" {
			totalStatements++
		}
	}

	return totalStatements > 0 && exportStatements == totalStatements
}

func findUseEffectCalls(node *sitter.Node, source []byte) {
	if node == nil {
		return
	}

	if node.Type() == "call_expression" {
		funcNode := node.ChildByFieldName("function")
		if funcNode != nil && funcNode.Type() == "identifier" && funcNode.Content(source) == "useEffect" {
			fmt.Println("found use Effect")
			argsNode := node.ChildByFieldName("arguments")
			if argsNode != nil && argsNode.NamedChildCount() > 0 {
				effectFunc := argsNode.NamedChild(0)

				if effectFunc != nil && (effectFunc.Type() == "arrow_function" || effectFunc.Type() == "function") {
					if effectFunc.ChildByFieldName("async") != nil {
					}

					body := effectFunc.ChildByFieldName("body")
					if body != nil && containsSetState(body, source) {
						if argsNode.NamedChildCount() == 1 {
							fmt.Printf("useEffect sets state but has no dependency array at")
						} else if argsNode.NamedChild(1).Type() == "array" && argsNode.NamedChild(1).ChildCount() == 0 {
							fmt.Printf("useEffect sets state with empty dep")
						}
					}
				}
			}
		}
	}

	for i := range int(node.ChildCount()) {
		findUseEffectCalls(node.Child(i), source)
	}
}

func containsSetState(node *sitter.Node, source []byte) bool {
	if node == nil {
		return false
	}

	if node.Type() == "call_expression" {
		funcNode := node.ChildByFieldName("function")
		if funcNode != nil && funcNode.Type() == "identifier" {
			name := funcNode.Content(source)
			if strings.HasPrefix(name, "set") {
				return true
			}
		}
	}

	for i := range int(node.ChildCount()) {
		if containsSetState(node.Child(i), source) {
			return true
		}
	}
	return false
}
