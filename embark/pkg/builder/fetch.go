package builder

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/containerum/containerum/embark/pkg/cgraph"
	"github.com/containerum/containerum/embark/pkg/models/requirements"
)

// Fetches requirements recursively to dir and building dependency graph
func (client *Client) FetchAllDeps(rootRequirements requirements.Requirements, dir string) (cgraph.SGraph, error) {
	//	if err := client.DownloadRequirements(dir, rootRequirements); err != nil {
	//		return err
	//	}
	var deps = requirements.NewQueue(len(rootRequirements.Dependencies))
	deps.Push(rootRequirements.Dependencies...)
	var getter = &http.Client{
		Timeout: 60 * time.Second,
	}
	var downloaded = map[string]bool{}

	var graph = make(cgraph.SGraph)
	graph.AddNode(Containerum, rootRequirements.Names()...)
	for dep := range deps {
		dep := dep

		fmt.Printf("Resolving %q, %d deps left\n", dep, len(deps))
		var depDep []string
		if !downloaded[dep.Name] {
			var depDir = path.Join(dir, dep.Name)
			if err := client.downloadDependency(getter, dir, dep); err != nil {
				fmt.Println(err)
				continue
			}
			var depReq, errDepReq = client.getRequirements(depDir)
			if errDepReq != nil {
				if !strings.Contains(errDepReq.Error(), ".yaml not found") {
					return nil, errDepReq
				}
			}

			var depChart, errLoadChart = client.LoadChartFromDir(depDir)
			if errLoadChart != nil {
				fmt.Println(errLoadChart)
				continue
			}

			if len(depChart.GetDependencies()) == 0 {
				fmt.Printf("\t%q depends on %v\n", dep.Name, depReq.Dependencies)
				deps.Push(depReq.Dependencies...)
				depDep = depReq.Names()
			} else {
				fmt.Printf("\tDeps of %q are already vendored in 'charts' dir\n", dep.Name)
			}
			downloaded[dep.Name] = true
		} else {
			fmt.Printf("\t%q is already fetched", dep)
		}

		fmt.Printf("\tAdding %q to graph\n", dep)
		graph.AddNode(dep.Name, depDep...)

		if len(deps) == 0 {
			close(deps)
		}
	}
	fmt.Printf("Charts to InstallChartWithDependencies:\n\t%s\n", strings.Join(graph.Nodes(), "\n\t"))
	return graph, nil
}
