package elee


type children []*node

type node struct{
	kind uint8 // 路由类型 0 静态路由 1 带参数路由 2 全匹配路由
	label byte // prefix 的第一个字符，根据label和kind来查询子节点
	parent *node	// 父节点
	staticChildrens children // 自己点列表
	methodHandler *methodHandler // 对应的handler
	params map[string]string // 参数
	prefix string // 前缀
}


