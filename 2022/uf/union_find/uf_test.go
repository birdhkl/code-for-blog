package union_find_test

import (
	"testing"
	"uf/union_find"
)

func TestUnionFind(t *testing.T) {
	type TestCase struct {
		data  [][2]int32 // [(自己所属序号，集合序号)]
		resut [][]int32  // 最终结果，按集合序号排列
		uf    union_find.UnionFind
	}

	testCases := []TestCase{}

	testCases = append(testCases, TestCase{
		data: [][2]int32{
			{0, 0},
			{1, 1},
			{2, 0},
			{3, 0},
			{4, 1},
			{5, 1},
			{6, 0},
			{7, 7},
		},
		resut: [][]int32{
			{0, 2, 3, 6},
			{1, 4, 5},
			{7},
		},
		uf: union_find.NewQuickFindUnionFind(10),
	})

	testCases = append(testCases, TestCase{
		data: [][2]int32{
			{0, 0},
			{1, 1},
			{2, 0},
			{3, 0},
			{4, 1},
			{5, 1},
			{6, 0},
			{7, 7},
		},
		resut: [][]int32{
			{0, 2, 3, 6},
			{1, 4, 5},
			{7},
		},
		uf: union_find.NewQuickUnionUnionFind(10),
	})

	testCases = append(testCases, TestCase{
		data: [][2]int32{
			{0, 0},
			{1, 1},
			{2, 0},
			{3, 0},
			{4, 1},
			{5, 1},
			{6, 0},
			{7, 7},
		},
		resut: [][]int32{
			{0, 2, 3, 6},
			{1, 4, 5},
			{7},
		},
		uf: union_find.NewWeightQuickUnionUnionFind(10),
	})

	testCases = append(testCases, TestCase{
		data: [][2]int32{
			{0, 0},
			{1, 1},
			{2, 0},
			{3, 0},
			{4, 1},
			{5, 1},
			{6, 0},
			{7, 7},
		},
		resut: [][]int32{
			{0, 2, 3, 6},
			{1, 4, 5},
			{7},
		},
		uf: union_find.NewCompressWeightQuickUnionUnionFind(10),
	})

	for caseNo, testCase := range testCases {
		for _, item := range testCase.data {
			testCase.uf.Add(item[0])
			testCase.uf.Union(item[1], item[0])
		}

		if int32(len(testCase.resut)) != testCase.uf.Count() {
			t.Errorf("case %d, set number wrong", caseNo)
		}

		if len(testCase.resut) != 1 {
			for setNo := 0; setNo < len(testCase.resut); setNo++ {
				for elementNo := 1; elementNo < len(testCase.resut[setNo]); elementNo++ {
					if !testCase.uf.Connected(testCase.resut[setNo][0], testCase.resut[setNo][elementNo]) {
						t.Errorf("case %d, set %d, not same set", caseNo, setNo)
					}
				}
				for otherCase := setNo + 1; otherCase < len(testCase.resut); otherCase++ {
					if testCase.uf.Connected(testCase.resut[setNo][0], testCase.resut[otherCase][0]) {
						t.Errorf("case %d, adjacent to case %d", caseNo, otherCase)
					}
				}
			}
		}
	}
}
