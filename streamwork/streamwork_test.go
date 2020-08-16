package streamwork

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triasteam/go-streamnet/types"
)

var (
	A types.Hash
	B types.Hash
	C types.Hash
	D types.Hash
	E types.Hash
	X types.Hash
	Y types.Hash
)

func (d *Dag) InitDag() {
	d.graph = make(map[types.Hash][]types.Hash)
	d.revGraph = make(map[types.Hash][]types.Hash)
	d.parentGraph = make(map[types.Hash]types.Hash)
	d.revParentGraph = make(map[types.Hash][]types.Hash)

	d.parentScore = make(map[types.Hash]float64)

	A := types.NewHashString("A")
	B := types.NewHashString("B")
	C := types.NewHashString("C")
	D := types.NewHashString("D")
	E := types.NewHashString("E")
	X := types.NewHashString("X")
	Y := types.NewHashString("Y")

	d.revGraph[A] = []types.Hash{B, C}
	d.revGraph[B] = []types.Hash{D, E}
	d.revGraph[C] = []types.Hash{D, E}
	d.revGraph[X] = []types.Hash{A}

	d.graph[D] = []types.Hash{B, C}
	d.graph[E] = []types.Hash{B, C}
	d.graph[C] = []types.Hash{A}
	d.graph[B] = []types.Hash{A}
	d.graph[A] = []types.Hash{X, Y}

	d.revParentGraph[A] = []types.Hash{B, C}
	d.revParentGraph[B] = []types.Hash{D, E}

	d.parentGraph[D] = B
	d.parentGraph[B] = A
	d.parentGraph[C] = A
	d.parentGraph[E] = B
	d.parentGraph[A] = X

	d.parentScore[A] = float64(3)
	d.parentScore[B] = float64(2)
	d.parentScore[C] = float64(2)
	d.parentScore[D] = float64(1)
	d.parentScore[E] = float64(1)

	d.genesis = A
}

func TestIfCovered(t *testing.T) {
	dag := &Dag{}
	dag.InitDag()
	cover := dag.IfCovered(B, A, []types.Hash{})
	a := assert.New(t)
	a.Equal(cover, false)
	// noCover := dag.IfCovered(D, A, []string{})
	// a.Equal(noCover, true)
}

func TestDiffSet(t *testing.T) {
	dag := &Dag{}
	dag.InitDag()
	children := dag.DiffSet(D, B, []types.Hash{})
	a := assert.New(t)
	a.Equal(2, len(children))
}

func TestStringToFloat(t *testing.T) {
	a := StringToFloat(types.NewHashString("nihao"))
	t.Log(a)
	b := StringToFloat(types.NewHashString("nihAo"))
	t.Log(b)
}

func TestGetPivot(t *testing.T) {
	dag := &Dag{}
	dag.InitDag()
	pivotHash := dag.getPivot(A)
	t.Log(pivotHash)
	a := assert.New(t)
	a.Equal(D, pivotHash)
}

func TestRemoveALl(t *testing.T) {
	v_1 := types.NewHashString("1")
	v_2 := types.NewHashString("2")
	v_3 := types.NewHashString("3")
	v_4 := types.NewHashString("4")
	v_5 := types.NewHashString("5")

	obj := []types.Hash{v_1, v_2, v_3, v_4, v_5}
	ano := []types.Hash{v_1, v_2, v_3, v_4, v_5}
	obj1 := RemoveAll(obj, ano)
	a := assert.New(t)
	a.Equal(0, len(obj1))

	ano = []types.Hash{v_1, v_2, v_4, v_5}
	obj2 := RemoveAll(obj, ano)
	a = assert.New(t)
	a.Equal(1, len(obj2))
	a.Equal(v_3, obj2[0])

	obj = []types.Hash{}
	ano = []types.Hash{v_1, v_2, v_4, v_5}
	obj2 = RemoveAll(obj, ano)
	a = assert.New(t)
	a.Equal(0, len(obj2))
}

func TestStreamWork(t *testing.T) {
	dag := &Dag{}
	dag.InitDag()
	totalOrder := dag.StreamWork(D)
	t.Log(totalOrder)
}

func TestSortByLvl(t *testing.T) {
	a := []types.Hash{B, C, A, D}
	b := make(map[types.Hash]int64)
	b[A] = 3
	b[B] = 2
	b[C] = 2
	b[D] = 1
	a = SortByLvl(a, b)
	t.Log(a)
}
