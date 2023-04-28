package apirouter

import (
	"bytes"
	"fmt"
	"moratal/apirouter/util"
	"sort"
)

var supportedParamTypes = map[string]struct{}{
	"int": {}, "string": {}, "bool": {}, "file": {},
}
var pathParamTypePriority = map[string]int{
	"bool":   0,
	"file":   1,
	"int":    2,
	"string": 3,
}

type PathParam struct {
	key string
	typ string
}

type segment struct {
	pattern   string     // 非叶子节点为'', 叶子节点为完整的path，如：'/api/nodes/v1/{id:int}'
	part      string     // 路由中的一部分，例如：pattern='/api/nodes' 时, part='nodes'
	children  []*segment // 子节点，例如 [nodes, audittype-logs, users]
	pathParam *PathParam // 动态路由'/api/nodes/{id:int}'的{id:int}, 存在类似{id:int}时，pathParam != nil, 否则为nil
}

/*
insert

	@param pattern string '/api/nodes/{id:int}'
	@param parts []string ["api", "nodes", "{id:int}"]
*/
func (s *segment) insert(method, pattern string, parts []string, height int) {
	// len(parts) == height: 说明已处理完所有节点，设置完毕后即可退出
	if len(parts) == height {
		s.pattern = pattern
		return
	}
	//当前节点为非叶子结点，则先查找是否存在匹配pattern节点，找到了则基于该child递归insert，
	//没找到则说明是一个新的path分支，新建child，并且继续递归insert
	part := parts[height]
	child := s.matchChild(part)
	if child == nil {
		child = &segment{part: part}
		// 若匹配到为动态路由，则设置动态参数
		key, typ := getDynamicPathParamInfo(part)
		if key != "" && typ != "" {
			if _, ok := supportedParamTypes[typ]; !ok {
				panic(fmt.Sprintf("unexpected param type '%s' in route pattern '%s %s'", typ, method, pattern))
			}
			child.pathParam = &PathParam{key: key, typ: typ}
		}
		s.children = append(s.children, child)
	}
	child.insert(method, pattern, parts, height+1)

}
func (s *segment) matchChild(part string) *segment {
	for _, child := range s.children {
		key, typ := getDynamicPathParamInfo(part)
		if child.part == part || (child.pathParam != nil && child.pathParam.key == key && child.pathParam.typ == typ) {
			return child
		}
	}
	return nil
}

// sortChildren
//
//	排序：优先做精确匹配
//	pattern1='[GET] /api/nodes/v1/{id:int}'
//	pattern2='[GET] /api/nodes/v1/1234'
//	e.g. 为保证reqInfo='[GET] /api/nodes/v1/1234' 优先匹配到 pattern2，需要保证nodes中的node2{pattern: pattern2}要先于node1{pattern: pattern1}
//	【注意】：目前的路由匹配无法很好地处理前缀均相同，仅path尾部的动态参数不同的path，例如：'/api/nodes/v1/list/{node_id:int}'和'/api/nodes/v1/list/{name:string}'
//	这会导致api处理时存在混乱，并且经过试验发现，iris也存在该类情况，仅能处理部分api-path。因此，在设计api时应避免该类情况发生。
func (s *segment) sortChildren() {
	if len(s.children) == 0 {
		return
	}
	nodes := s.children
	sort.Slice(nodes, func(i, j int) bool {
		//路由排序
		switch {
		// 1. 无动态路由参数的pattern, e.g. '/api/nodes/v1/list'
		case nodes[i].pathParam == nil:
			return true
			// 2. 存在动态路由参数的pattern, e.g. '/api/nodes/v1/{node_id:int}
		case nodes[i].pathParam != nil && nodes[j].pathParam == nil:
			return false
		case nodes[i].pathParam != nil && nodes[j].pathParam != nil:
			return pathParamTypePriority[nodes[i].pathParam.typ] <= pathParamTypePriority[nodes[j].pathParam.typ]
		default:
			return false
		}
	})
	s.children = nodes
	for i := 0; i < len(s.children); i++ {
		s.children[i].sortChildren()
	}
}

// getDynamicPathParamInfo 支持类似'{id:int}'格式的path-param, '{key:type}'中':'两边不允许有空格，否则会导致解析为空
func getDynamicPathParamInfo(segment string) (key, typ string) {
	part := util.String2Bytes(segment)
	if part[0] == '{' && part[len(part)-1] == '}' {
		indexByte := bytes.IndexByte(part, ':')
		if indexByte > 0 {
			return string(part[1:indexByte]), string(part[indexByte+1:])
		}
		return string(part), "string"
	}
	return "", ""
}
