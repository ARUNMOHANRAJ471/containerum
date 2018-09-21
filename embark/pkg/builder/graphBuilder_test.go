package builder

import (
	"os"
	"path"
	"testing"

	"github.com/containerum/containerum/embark/pkg/models/components"
)

func TestBuildGraph(test *testing.T) {
	var testdir = path.Join("./testdata", "buildGraph") //path.Join(os.TempDir(), "embark", "testBuildGraph")
	os.MkdirAll(testdir, os.ModeDir|os.ModePerm)
	var cont = components.Components{
		"mongodb": components.Component{
			Repo:    "https://charts.containerum.io",
			Version: "3.0.4",
			Values:  map[string]interface{}{},
			Objects: []string{
				"configmap",
				"svc-standalone",
			},
		},
	}
	var downloadErr = DowloadComponents(testdir, cont)
	if downloadErr != nil {
		test.Fatal(downloadErr)
	}

	var rendered, renderErr = RenderComponents(testdir, cont)
	if renderErr != nil {
		test.Fatal(renderErr)
	}
	for _, obj := range rendered {
		for name, obj := range obj.Objects {
			test.Log("\n", name, "-> ", obj)
		}
	}
	var _, buildGraphErr = BuildGraph(testdir, rendered)
	if buildGraphErr != nil {
		test.Fatal(buildGraphErr)
	}
}
