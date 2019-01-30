/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : stack.go
#   Created       : 2019-01-30 11:35:06
#   Describe      :
#
# ====================================================*/
package stack

// Item 栈中的每个元素
type Item struct {
	Value []byte
}

// Stack 栈
type Stack struct {
	items []*Item
}

// New 初始化一个栈
func New() *Stack {
	return &Stack{}
}

// Push 入栈
func (s *Stack) Push(t *Item) {
	s.items = append(s.items, t)
}

// Pop 出栈
func (s *Stack) Pop() *Item {
	if len(s.items) == 0 {
		return nil
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

// Value 获取栈中所有值
func (s *Stack) Value() []*Item {
	return s.items
}

// Flush 清空栈
func (s *Stack) Flush() {
	s.items = []*Item{}
}
