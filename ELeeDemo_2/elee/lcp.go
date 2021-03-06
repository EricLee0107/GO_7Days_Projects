package elee

import (
	"fmt"
	"net/http"
	"strings"
)

type children []*node

type node struct{
	kind uint8 // 路由类型 0 静态路由 1 带参数路由 2 全匹配路由
	label byte // prefix 的第一个字符，根据label和kind来查询子节点
	parent *node	// 父节点
	staticChildrens children // 子节点列表
	methodHandler *methodHandler // 对应的handler
	pnames          []string // 路径参数（只有当kind为1或者2是才有）
	prefix string // 前缀
	ppath string
	elee *Elee
	paramChildren   *node
	anyChildren     *node

}




var (
	methods = [...]string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace,
	}
)



const (
	skind uint8 = iota
	pkind
	akind
)



// /index/test
// /index/user
// /index/user/elee
// /index/:test
// /index/*

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
// 核心函数，构建字典树
func (cn *node) insert(method, path string, h HandlerFunc, t uint8, ppath string, pnames []string) {
	// 调整最大参数
	l := len(pnames)
	if *cn.elee.maxParam < l {
		*cn.elee.maxParam = l
	}

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
			if cn.paramChildren != nil {
				cn.paramChildren.parent = n
			}
			if cn.anyChildren != nil {
				cn.anyChildren.parent = n
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
			cn.paramChildren = nil
			cn.anyChildren = nil

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
				cn.methodHandler.addHandler(method, h)
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
		paramChildren:   paramChildren,
		anyChildren:     anyChildren,
	}
}


func (n *node) addStaticChild(c *node) {
	n.staticChildrens = append(n.staticChildrens, c)
}
func (n *node) findChild(l byte) *node {
	for _, c := range n.staticChildrens {
		if c.label == l {
			return c
		}
	}
	return nil
}

func (n *node) findChildWithLabel(l byte) *node {
	for _, c := range n.staticChildrens {
		if c.label == l {
			return c
		}
	}
	if l == byte(':') {
		return n.paramChildren
	}
	if l == byte('*') {
		return n.anyChildren
	}
	return nil
}

func (n *node) findStaticChild(l byte) *node {
	for _, c := range n.staticChildrens {
		if c.label == l {
			return c
		}
	}
	return nil
}


func (n *node) checkMethodNotAllowed() HandlerFunc {
	for _, m := range methods {
		if h := n.methodHandler.findHandler(m); h != nil {
			return func(c Context) error {
		return fmt.Errorf("error1:%d",404)
	}
		}
	}
	return func(c Context) error {
		return fmt.Errorf("error2:%d",404)
	}
}


// 通过method和path查找住的的handler，解析URL参数并把参数放入context
func (cn *node) Find(method, path string, c Context) {
	ctx := c.(*context)
	ctx.path = path

	var (
		search  = path
		child   *node         // 子节点
		n       int           // 参数计数器
		nk      uint8         // 下一个节点的kind
		nn      *node         // 下一个节点
		ns      string        // 下一个search字串
		pvalues = ctx.pvalues // Use the internal slice so the interface can keep the illusion of a dynamic slice
	)

	// 搜索顺序 static > param > any
	for {
		if search == "" {
			break
		}

		pl := 0 // Prefix length
		l := 0  // LCP length

		if cn.label != ':' {
			sl := len(search)
			pl = len(cn.prefix)

			// LCP
			max := pl
			if sl < max {
				max = sl
			}
			// 找到共同前缀的起始点
			for ; l < max && search[l] == cn.prefix[l]; l++ {
			}
		}

		if l == pl {
			// 重合，继续搜索
			search = search[l:]
			// Finish routing if no remaining search and we are on an leaf node
			if search == "" && (nn == nil || cn.parent == nil || cn.ppath != "") {
				break
			}
			// Handle special case of trailing slash route with existing any route (see #1526)
			if search == "" && path[len(path)-1] == '/' && cn.anyChildren != nil {
				goto Any
			}
		}

		// Attempt to go back up the tree on no matching prefix or no remaining search
		if l != pl || search == "" {
			if nn == nil { // Issue #1348
				return // Not found
			}
			cn = nn
			search = ns
			if nk == pkind {
				goto Param
			} else if nk == akind {
				goto Any
			}
		}

		// Static 节点
		if child = cn.findStaticChild(search[0]); child != nil {
			// Save next
			if cn.prefix[len(cn.prefix)-1] == '/' { // Issue #623
				nk = pkind
				nn = cn
				ns = search
			}
			cn = child
			continue
		}

	Param:
		// Param 节点
		if child = cn.paramChildren; child != nil {
			// Issue #378
			if len(pvalues) == n {
				continue
			}

			// Save next
			if cn.prefix[len(cn.prefix)-1] == '/' { // Issue #623
				nk = akind
				nn = cn
				ns = search
			}

			cn = child
			i, l := 0, len(search)
			for ; i < l && search[i] != '/'; i++ {
			}
			pvalues[n] = search[:i]
			n++
			search = search[i:]
			continue
		}

	Any:
		// Any 节点
		if cn = cn.anyChildren; cn != nil {
			// If any node is found, use remaining path for pvalues
			pvalues[len(cn.pnames)-1] = search
			break
		}

		// No node found, continue at stored next node
		// or find nearest "any" route
		if nn != nil {
			// No next node to go down in routing (issue #954)
			// Find nearest "any" route going up the routing tree
			search = ns
			np := nn.parent
			// Consider param route one level up only
			if cn = nn.paramChildren; cn != nil {
				pos := strings.IndexByte(ns, '/')
				if pos == -1 {
					// If no slash is remaining in search string set param value
					if len(cn.pnames) > 0 {
						pvalues[len(cn.pnames)-1] = search
					}
					break
				} else if pos > 0 {
					// Otherwise continue route processing with restored next node
					cn = nn
					nn = nil
					ns = ""
					goto Param
				}
			}
			// No param route found, try to resolve nearest any route
			for {
				np = nn.parent
				if cn = nn.anyChildren; cn != nil {
					break
				}
				if np == nil {
					break // no further parent nodes in tree, abort
				}
				var str strings.Builder
				str.WriteString(nn.prefix)
				str.WriteString(search)
				search = str.String()
				nn = np
			}
			if cn != nil { // use the found "any" route and update path
				pvalues[len(cn.pnames)-1] = search
				break
			}
		}
		return // Not found

	}

	ctx.handler = cn.methodHandler.findHandler(method)
	ctx.path = cn.ppath
	ctx.pnames = cn.pnames

	// NOTE: Slow zone...
	if ctx.handler == nil {
		ctx.handler = cn.checkMethodNotAllowed()

		// Dig further for any, might have an empty value for *, e.g.
		// serving a directory. Issue #207.
		if cn = cn.anyChildren; cn == nil {
			return
		}
		if h := cn.methodHandler.findHandler(method); h != nil {
			ctx.handler = h
		} else {
			ctx.handler = cn.checkMethodNotAllowed()
		}
		ctx.path = cn.ppath
		ctx.pnames = cn.pnames
		pvalues[len(cn.pnames)-1] = ""
	}

	return
}

