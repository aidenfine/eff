package analyzer

import (
	"context"
	"fmt"
	"os"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

var useEffectCount int
var importNodes []string

func AnalyzeFile(path string) {
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
	// find dead code
	findImportNodes(root, source, path)
	// fmt.Println(path)
	fmt.Println(useEffectCount, "useEffect count")

}
func findImportNodes(node *sitter.Node, source []byte, path string) {
	if node == nil {
		return
	}
	if node.Type() == "import_statement" {
		importNodes = append(importNodes, node.Content(source))
	}
	for i := range int(node.ChildCount()) {
		child := node.Child(i)
		findImportNodes(child, source, path)

	}

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
			}
		}
	}

	for i := range int(node.ChildCount()) {
		findUseEffectCalls(node.Child(i), source, path)
	}
}
