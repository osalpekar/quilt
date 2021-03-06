package inspect

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/quilt/quilt/util"
)

func stripExtension(configPath string) string {
	ext := filepath.Ext(configPath)
	return strings.TrimSuffix(configPath, ext)
}

func viz(configPath string, graph Graph, outputFormat string) {
	slug := stripExtension(configPath)
	dot := makeGraphviz(graph)
	graphviz(outputFormat, slug, dot)
}

func makeGraphviz(graph Graph) string {
	dotfile := "strict digraph {\n"

	var lines []string
	for _, edge := range graph.GetConnections() {
		lines = append(lines,
			fmt.Sprintf(
				"    %q -> %q;\n",
				edge.From,
				edge.To,
			),
		)
	}

	sort.Strings(lines)
	for _, line := range lines {
		dotfile += line + "\n"
	}

	dotfile += "}\n"

	return dotfile
}

// Graphviz generates a specification for the graphviz program that visualizes the
// communication graph of a stitch.
func graphviz(outputFormat string, slug string, dot string) {
	f, err := util.AppFs.Create(slug + ".dot")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write([]byte(dot))
	if outputFormat == "graphviz" {
		return
	}
	defer exec.Command("rm", slug+".dot").Run()

	// Dependencies:
	// - easy-graph (install Graph::Easy from cpan)
	// - graphviz (install from your favorite package manager)
	var writeGraph *exec.Cmd
	switch outputFormat {
	case "ascii":
		writeGraph = exec.Command("graph-easy", "--input="+slug+".dot",
			"--as_ascii")
	case "pdf":
		writeGraph = exec.Command("dot", "-Tpdf", "-o", slug+".pdf",
			slug+".dot")
	}
	writeGraph.Stdout = os.Stdout
	writeGraph.Stderr = os.Stderr
	writeGraph.Run()
}
