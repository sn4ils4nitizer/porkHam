package utils

import (
	"os"
	"path/filepath"
)

type TreeNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	Children []*TreeNode `json:"children,omitempty"`
	IsDir    bool        `json:"isDir"`
}

func BuildTree(rootPath string) (*TreeNode, error) {
	rootInfo, err := os.Stat(rootPath)
	if err != nil {
		return nil, err
	}

	rootNode := &TreeNode{
		Name:  rootInfo.Name(),
		Path:  rootPath,
		IsDir: rootInfo.IsDir(),
	}

	if rootNode.IsDir {
		err = buildTreeRecursive(rootPath, rootNode)
		if err != nil {
			return nil, err
		}
	}
	return rootNode, nil
}

func buildTreeRecursive(path string, parent *TreeNode) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())
		childNode := &TreeNode{
			Name:  entry.Name(),
			Path:  childPath,
			IsDir: entry.IsDir(),
		}

		if entry.IsDir() {
			err := buildTreeRecursive(childPath, childNode)
			if err != nil {
				return err
			}
		}

		parent.Children = append(parent.Children, childNode)
	}
	return nil
}
