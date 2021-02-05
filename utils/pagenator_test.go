package utils

import (
	"testing"

	"github.com/triasteam/go-streamnet/types"
)

func TestPaging(t *testing.T) {
	data := make([]types.Hash, 0)
	hash1 := types.NewHashString("1")
	data = append(data, hash1)
	hash2 := types.NewHashString("2")
	data = append(data, hash2)
	hash3 := types.NewHashString("3")
	data = append(data, hash3)
	hash4 := types.NewHashString("4")
	data = append(data, hash4)
	hash5 := types.NewHashString("5")
	data = append(data, hash5)
	hash6 := types.NewHashString("6")
	data = append(data, hash6)

	newArray := Paging(data, 1, 5)
	t.Logf("new array is : %v \n", newArray)

	newArray1 := Paging(data, 1, 10)
	t.Logf("new array is : %v \n", newArray1)

	newArray2 := Paging(data, -1, 5)
	t.Logf("new array is : %v \n", newArray2)

	newArray3 := Paging(data, 2, 5)
	t.Logf("new array is : %v \n", newArray3)

	// 异常情况，data为空
	// data = make([]types.Hash, 0)
	// data = append(data)
	// res := Paging(data, -1, 5)
	// t.Logf("%v", res)
}
