package union_find

type QuickUnionUnionFind struct {
	ids   []int32
	count int32
}

func NewQuickUnionUnionFind(cap int32) UnionFind {
	return &QuickUnionUnionFind{
		ids:   make([]int32, 0, cap),
		count: 0,
	}
}

func (uf *QuickUnionUnionFind) Add(p int32) {
	if p < int32(len(uf.ids)) {
		return
	}
	for p >= int32(len(uf.ids)) {
		uf.ids = append(uf.ids, int32(len(uf.ids)))
		uf.count += 1
	}
}

func (uf *QuickUnionUnionFind) Count() int32 {
	return uf.count
}

func (uf *QuickUnionUnionFind) Find(p int32) int32 {
	for p != uf.ids[p] {
		p = uf.ids[p]
	}
	return p
}

func (uf *QuickUnionUnionFind) Union(p int32, q int32) {
	pi := uf.Find(p)
	qi := uf.Find(q)
	if pi == qi {
		return
	}
	uf.ids[qi] = pi
	uf.count -= 1
}

func (uf *QuickUnionUnionFind) Connected(p int32, q int32) bool {
	return uf.Find(p) == uf.Find(q)
}
