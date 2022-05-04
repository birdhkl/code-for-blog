package union_find

type WeightQuickUnionUnionFind struct {
	ids     []int32
	weights []int32
	count   int32
}

func NewWeightQuickUnionUnionFind(cap int32) UnionFind {
	return &WeightQuickUnionUnionFind{
		ids:     make([]int32, 0, cap),
		weights: make([]int32, 0, cap),
		count:   0,
	}
}

func (uf *WeightQuickUnionUnionFind) Add(p int32) {
	if p < int32(len(uf.ids)) {
		return
	}
	for p >= int32(len(uf.ids)) {
		uf.ids = append(uf.ids, int32(len(uf.ids)))
		uf.weights = append(uf.weights, 1)
		uf.count += 1
	}
}

func (uf *WeightQuickUnionUnionFind) Find(p int32) int32 {
	for p != uf.ids[p] {
		p = uf.ids[p]
	}
	return p
}

func (uf *WeightQuickUnionUnionFind) Connected(p int32, q int32) bool {
	return uf.Find(p) == uf.Find(q)
}

func (uf *WeightQuickUnionUnionFind) Union(p int32, q int32) {
	pi := uf.Find(p)
	qi := uf.Find(q)
	if pi == qi {
		return
	}
	// p所属的树是一颗更大的树
	if uf.weights[pi] >= uf.weights[qi] {
		uf.ids[qi] = pi
		uf.weights[pi] = uf.weights[qi] + uf.weights[pi]
		uf.weights[qi] = 0
	} else {
		// q所属的树是一颗更大的树
		uf.ids[pi] = qi
		uf.weights[qi] = uf.weights[qi] + uf.weights[pi]
		uf.weights[pi] = 0
	}
	uf.count -= 1
}

func (uf *WeightQuickUnionUnionFind) Count() int32 {
	return uf.count
}
