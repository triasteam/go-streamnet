package utils

import "github.com/triasteam/go-streamnet/types"

// Paging ...
func Paging(data []types.Hash, currentPage, pageSize int) []types.Hash {
	if len(data) < 0 {
		return nil
	}
	totalCount := len(data)
	m := Mod(totalCount, pageSize)

	tp := totalCount / pageSize

	totalPage := 0
	if m == 0 {
		totalPage = tp
	} else {
		totalPage = tp + 1
	}
	if currentPage == -1 {
		currentPage = totalPage
	}

	startIndex := (currentPage - 1) * pageSize
	if startIndex > totalCount {
		return nil
	}

	endIndex := currentPage * pageSize
	if endIndex > totalCount {
		endIndex = totalCount
	}

	lastPageSize := len(data[startIndex:endIndex])

	newArray := make([]types.Hash, lastPageSize)

	copy(newArray[:], data[startIndex:endIndex])

	return newArray
}
