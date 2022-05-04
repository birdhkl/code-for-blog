package union_find

type CompressWeightQuickUnionUnionFind struct {
	WeightQuickUnionUnionFind
}

func NewCompressWeightQuickUnionUnionFind(cap int) UnionFind {
	return &CompressWeightQuickUnionUnionFind{
		WeightQuickUnionUnionFind: WeightQuickUnionUnionFind{
			ids:     make([]int32, 0, cap),
			weights: make([]int32, 0, cap),
			count:   0,
		},
	}
}

func (uf *CompressWeightQuickUnionUnionFind) Find(p int32) int32 {
	// 对半压缩路径
	for p != uf.ids[p] {
		uf.ids[p] = uf.ids[uf.ids[p]]
		p = uf.ids[p]
	}
	return p
}
