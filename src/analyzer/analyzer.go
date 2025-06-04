package analyzer

// what we could also do here is keep track of a score and weight
// for example if file name is index.ts, index.tsx, ...  higher chance of barrel import file increase weight

import (
	"context"
	"fmt"
	"os"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

var useEffectCount int
var exportNodes []string
var lineCount int

func AnalyzeFile(path string) bool {
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to read %s: %v\n", path, err)
		return false
	}
	parser := sitter.NewParser()
	parser.SetLanguage(typescript.GetLanguage())
	tree, _ := parser.ParseCtx(context.Background(), nil, source)
	root := tree.RootNode()

	lineCount = int(root.ChildCount())

	findUseEffectCalls(root, source, path)
	if isBarrelFile(root, source) {
		fmt.Println("is barrel export")
		return true
	}
	return false

}
func isBarrelFile(root *sitter.Node, source []byte) bool {
	exportStatements := 0
	totalStatements := 0

	for i := 0; i < int(root.NamedChildCount()); i++ {
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
						// startByte := node.StartByte()
						// endByte := node.EndByte()
						// useEffectCall := source[startByte:endByte]
					}
				}
			}
		}
	}

	for i := range int(node.ChildCount()) {
		findUseEffectCalls(node.Child(i), source, path)
	}
}
