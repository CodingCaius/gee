package gee

import (
    "fmt"
    "strings"
)

type node struct {
    pattern string  //存储路径，如果不为空，说明当前节点是一个叶子节点
    part string  //存储的是路径的一段
    children []*node //子节点切片
    isWild bool  //是否是通配符
}

func (n *node) String() string {
    return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

//插入路径
//参数 pattern 是待插入的路径模式，参数 parts 是路径模式的各个部分切片，参数 height 表示当前插入的部分在切片中的位置
func (n *node) insert(pattern string, parts []string, height int) {
    if len(parts) == height {
        n.pattern = pattern
        return
    }

    part := parts[height]
    child := n.matchChild(part)
    if child == nil {
        child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
        n.children = append(n.children, child)
    }
    child.insert(pattern, parts, height + 1)
}

func (n *node) search(parts []string, height int) *node {
    if len(parts) == height || strings.HasPrefix(n.part, "*") {
        if n.pattern == "" {
            return nil
        }
        return n
    }

    part := parts[height]
    //获取所有与当前部分匹配的子节点
    children := n.matchChildren(part)

    for _, child := range children {
        result := child.search(parts, height + 1)
        if result != nil {
            return result
        }
    }

    return nil
}

//遍历路由树中的节点，并将 具有模式的节点 添加到给定的切片列表中
func (n *node) travel(list *([]*node)) {
    if n.pattern != "" {
        *list = append(*list, n)
    }
    for _, child := range n.children {
        child.travel(list)
    }
}

//查找首个匹配的子节点
func (n *node) matchChild(part string) *node {
    for _, child := range n.children {
        if child.part == part || child.isWild {
            return child
        }
    }
    return nil
}

//查找所有匹配的子节点
func (n *node) matchChildren(part string) []*node {
    nodes := make([]*node, 0)
    for _, child := range n.children {
        if child.part == part || child.isWild {
            nodes = append(nodes, child)
        }
    }
    return nodes
}