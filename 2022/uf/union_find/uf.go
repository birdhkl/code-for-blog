package union_find

// UnionFind 并查集接口
type UnionFind interface {
	// 新增集合
	Add(p int32)
	// 查询集合编号，不存在则返回-1
	Find(p int32) int32
	// 集合是否相交
	Connected(p int32, q int32) bool
	// 合并集合
	Union(p int32, q int32)
	// 集合数量
	Count() int32
}
