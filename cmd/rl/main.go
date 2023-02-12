package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
)

var ignored_paths = map[string]bool{
	".git": true,
}

func main() {
	var paths []string
	if len(os.Args) == 1 {
		paths = []string{"."}
	} else {
		paths = os.Args[1:len(os.Args)]
	}
		
	for _, path := range paths {
		dfs(NewPath(path), println)
	}
}

func println(s string) {
	fmt.Println(s)
}

func dfs(root *path, action func(string)) {
	stack := []*path{NewPath()}

	for stackDepth := len(stack); stackDepth > 0; stackDepth = len(stack) {
		parent := stack[len(stack)-1]
		stack = stack[0:len(stack)-1]

		entries, err := ioutil.ReadDir(root.JoinPath(parent).String())
		if err != nil {
			log.Fatal(err)
		}

		dirs := []string{}
		files := []string{}

		for _, entry := range entries {
			if (entry.IsDir()) {
				if (!ignored_paths[entry.Name()]) {
					dirs = append(dirs, entry.Name())
				}
			} else {
				files = append(files, entry.Name())
			}
		}

		sort.Strings(files)

		for _, file := range files {
			if root.String() == "." {
				action(parent.Join(file).String())
			} else {
				action(root.JoinPath(parent.Join(file)).String())
			}
		}

		sort.Slice(dirs, func(i, j int) bool {
			return dirs[i] > dirs[j];
		})

		for _, dir := range dirs {
			stack = append(stack, parent.Join(dir))
		}
	}
}

type path struct {
	tokens []string
}

func NewPath(tokens ...string) *path {
	p := path{}
	p.tokens = tokens
	return &p
}

func (p *path) Join(tokens ...string) *path {
	joinedTokens := []string{}
	joinedTokens = append(joinedTokens, p.tokens...)
	joinedTokens = append(joinedTokens, tokens...)
	return NewPath(joinedTokens...)
}

func (p *path) JoinPath(p2 *path) *path {
	return p.Join(p2.tokens...)
}

func (p *path) String() string {
	return filepath.Join(p.tokens...)
}
