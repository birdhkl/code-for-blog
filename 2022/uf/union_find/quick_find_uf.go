package union_find

type QuickFindUnionFind struct {
	ids   []int32
	count int32
}

func NewQuickFindUnionFind(cap int32) UnionFind {
	return &QuickFindUnionFind{
		ids:   make([]int32, 0, cap),
		count: 0,
	}
}

func (uf *QuickFindUnionFind) Add(p int32) {
	if p < int32(len(uf.ids)) {
		return
	}
	for p >= int32(len(uf.ids)) {
		uf.ids = append(uf.ids, int32(len(uf.ids)))
		uf.count += 1
	}
}

func (uf *QuickFindUnionFind) Count() int32 {
	return uf.count
}

func (uf *QuickFindUnionFind) Find(p int32) int32 {
	return uf.ids[p]
}

func (uf *QuickFindUnionFind) Connected(p int32, q int32) bool {
	return uf.Find(p) == uf.Find(q)
}

func (uf *QuickFindUnionFind) Union(p int32, q int32) {
	pi := uf.Find(p)
	qi := uf.Find(q)
	if pi == qi {
		return
	}
	for i, id := range uf.ids {
		if id == qi {
			uf.ids[i] = pi
		}
	}
	uf.count -= 1
}
