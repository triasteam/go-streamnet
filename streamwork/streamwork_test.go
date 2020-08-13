package streamwork

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func (d *Dag) InitDag() {
	d.graph = make(map[string][]string)
	d.revGraph = make(map[string][]string)
	d.parentGraph = make(map[string]string)
	d.revParentGraph = make(map[string][]string)

	d.parentScore = make(map[string]float64)

	d.revGraph["A"] = []string{"B", "C"}
	d.revGraph["B"] = []string{"D", "E"}
	d.revGraph["C"] = []string{"D", "E"}
	d.revGraph["X"] = []string{"A"}

	d.graph["D"] = []string{"B", "C"}
	d.graph["E"] = []string{"B", "C"}
	d.graph["C"] = []string{"A"}
	d.graph["B"] = []string{"A"}
	d.graph["A"] = []string{"X", "Y"}

	d.revParentGraph["A"] = []string{"B", "C"}
	d.revParentGraph["B"] = []string{"D", "E"}

	d.parentGraph["D"] = "B"
	d.parentGraph["B"] = "A"
	d.parentGraph["C"] = "A"
	d.parentGraph["E"] = "B"
	d.parentGraph["A"] = "X"

	d.parentScore["A"] = float64(3)
	d.parentScore["B"] = float64(2)
	d.parentScore["C"] = float64(2)
	d.parentScore["D"] = float64(1)
	d.parentScore["E"] = float64(1)

	d.genesis = "A"
}

func TestIfCovered(t *testing.T) {
	dag := &Dag{}
	dag.InitDag()
	cover := dag.IfCovered("B", "A", []string{})
	a := assert.New(t)
	a.Equal(cover, false)
	// noCover := dag.IfCovered("D", "A", []string{})
	// a.Equal(noCover, true)
}

func TestDiffSet(t *testing.T) {
	dag := &Dag{}
	dag.InitDag()
	children := dag.DiffSet("D", "B", []string{})
	a := assert.New(t)
	a.Equal(2, len(children))
}

func TestStringToFloat(t *testing.T) {
	a := StringToFloat("nihao")
	t.Log(a)
	b := StringToFloat("nihAo")
	t.Log(b)
}

func TestGetPivot(t *testing.T) {
	dag := &Dag{}
	dag.InitDag()
	pivotHash := dag.getPivot("A")
	t.Log(pivotHash)
	a := assert.New(t)
	a.Equal("D", pivotHash)
}

func TestRemoveALl(t *testing.T) {
	obj := []string{"1", "2", "3", "4", "5"}
	ano := []string{"1", "2", "3", "4", "5"}
	obj1 := RemoveAll(obj, ano)
	a := assert.New(t)
	a.Equal(0, len(obj1))

	ano = []string{"1", "2", "4", "5"}
	obj2 := RemoveAll(obj, ano)
	a = assert.New(t)
	a.Equal(1, len(obj2))
	a.Equal("3", obj2[0])

	obj = []string{}
	ano = []string{"1", "2", "4", "5"}
	obj2 = RemoveAll(obj, ano)
	a = assert.New(t)
	a.Equal(0, len(obj2))
}

func TestStreamWork(t *testing.T) {
	dag := &Dag{}
	dag.InitDag()
	totalOrder := dag.StreamWork("D")
	t.Log(totalOrder)
}

func TestSortByLvl(t *testing.T) {
	a := []string{"B", "C", "A", "D"}
	b := make(map[string]int64)
	b["A"] = 3
	b["B"] = 2
	b["C"] = 2
	b["D"] = 1
	a = SortByLvl(a, b)
	t.Log(a)
}
