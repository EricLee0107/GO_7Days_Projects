package elee


type children []*node

type node struct{
	kind uint8 // 路由类型 0 静态路由 1 带参数路由 2 全匹配路由
	label byte // prefix 的第一个字符，根据label和kind来查询子节点
	parent *node	// 父节点
	staticChildrens children // 自己点列表
	methodHandler *methodHandler // 对应的handler
	pnames          []string // 路径参数（只有当kind为1或者2是才有）
	prefix string // 前缀
	ppath string
	elee *Elee
}

const (
	skind uint8 = iota
	pkind
	akind
)



// Add方法为路由器添加一个新的路径和对应的handler
func (n *node) Add(method, path string, h HandlerFunc) {
	// 验证路径合法性
	if path == "" {
		path = "/"
	}
	// 规范化路径
	if path[0] != '/' {
		path = "/" + path
	}
	pnames := []string{} // 路径参数
	ppath := path        // 原始路径

	for i, l := 0, len(path); i < l; i++ {
		// 参数路径
		if path[i] == ':' {
			j := i + 1

			n.insert(method, path[:i], nil, skind, "", nil)
			// 找到参数路径的参数
			for ; i < l && path[i] != '/'; i++ {
			}
			// 把参数路径存入pnames
			pnames = append(pnames, path[j:i])
			// 拼接路径，继续查找是否还有参数路径
			path = path[:j] + path[i:]
			i, l = j, len(path)
			// 已经结束，插入参数路径节点
			if i == l {
				n.insert(method, path[:i], h, pkind, ppath, pnames)
			} else {
				n.insert(method, path[:i], nil, pkind, "", nil)
			}
			// 全匹配路径
		} else if path[i] == '*' {
			n.insert(method, path[:i], nil, skind, "", nil)
			// 全匹配路径参数都是 *
			pnames = append(pnames, "*")
			n.insert(method, path[:i+1], h, akind, ppath, pnames)
		}
	}
	// 普通路径
	n.insert(method, path, h, skind, ppath, pnames)
}

// 核心函数，构建字典树
func (n *node) insert(method, path string, h HandlerFunc, t uint8, ppath string, pnames []string) {
	// 调整最大参数
	l := len(pnames)
	if *n.elee.maxParam < l {
		*n.elee.maxParam = l
	}

	cn := n.elee.tree // 当前节点root
	if cn == nil {
		panic("echo: invalid method")
	}
	search := path

	for {
		sl := len(search)
		pl := len(cn.prefix)
		l := 0

		// LCP
		max := pl
		if sl < max {
			max = sl
		}
		// 找到共同前缀的位置 例如users/ 和 users/new 的共同前缀为users/
		for ; l < max && search[l] == cn.prefix[l]; l++ {
		}

		if l == 0 {
			// root节点处理
			cn.label = search[0]
			cn.prefix = search
			if h != nil {
				cn.kind = t
				cn.methodHandler.addHandler(method, h)
				cn.ppath = ppath
				cn.pnames = pnames
			}
		} else if l < pl {

			// 分离共同前缀 users/和users/new 创建一个prefix为new的节点()
			n := newNode(cn.kind, cn.prefix[l:], cn, cn.staticChildrens, cn.methodHandler, cn.ppath, cn.pnames, cn.paramChildren, cn.anyChildren)

			// Update parent path for all children to new node
			// 将当前节点的所有子节点的父改为新的节点new
			for _, child := range cn.staticChildrens {
				child.parent = n
			}

			// Reset parent node
			cn.kind = skind
			cn.label = cn.prefix[0]
			cn.prefix = cn.prefix[:l]
			// 清空当前节点的所有子节点
			cn.staticChildrens = nil
			cn.methodHandler = new(methodHandler)
			cn.ppath = ""
			cn.pnames = nil

			// 将新创建的prefix为new的节点加到当前节点的子节点中
			cn.addStaticChild(n)

			if l == sl {
				// At parent node
				cn.kind = t
				cn.methodHandler.addHandler(method, h)
				cn.ppath = ppath
				cn.pnames = pnames
			} else {
				// Create child node
				n = newNode(t, search[l:], cn, nil, new(methodHandler), ppath, pnames, nil, nil)
				n.methodHandler.addHandler(method, h)
				// Only Static children could reach here
				cn.addStaticChild(n)
			}
		} else if l < sl {
			search = search[l:]
			// 找到lable一样的节点，用lable来判断共同前缀
			c := cn.findChildWithLabel(search[0])
			if c != nil {
				// 找到共同节点，继续
				cn = c
				continue
			}
			// 创建子节点
			n := newNode(t, search, cn, nil, new(methodHandler), ppath, pnames, nil, nil)
			n.methodHandler.addHandler(method, h)
			switch t {
			case skind:
				cn.addStaticChild(n)
			case pkind:
				cn.paramChildren = n
			case akind:
				cn.anyChildren = n
			}
		} else {
			// 节点已经存在
			if h != nil {
				cn.addHandler(method, h)
				cn.ppath = ppath
				if len(cn.pnames) == 0 { // Issue #729
					cn.pnames = pnames
				}
			}
		}
		return
	}
}



func newNode(t uint8, pre string, p *node, sc children, mh *methodHandler, ppath string, pnames []string, paramChildren, anyChildren *node) *node {
	return &node{
		kind:            t,
		label:           pre[0],
		prefix:          pre,
		parent:          p,
		staticChildrens: sc,
		ppath:           ppath,
		pnames:          pnames,
		methodHandler:   mh,
	}
}


