package tree

type BinaryTree struct {
	Root    *BinaryNode
	Min     *BinaryNode
	Compare func(lKey any, rKey any) int
}

type BinaryNode struct {
	Parent *BinaryNode
	Key    any
	Values []any
	Left   *BinaryNode
	Right  *BinaryNode
}

func (t *BinaryTree) Push(key any, value any) {
	if t.Root == nil {
		t.Root = &BinaryNode{
			Parent: nil,
			Key:    key,
			Values: []any{value},
			Left:   nil,
			Right:  nil,
		}

		t.Min = t.Root
		return
	}

	t.Root.Push(t, &BinaryNode{
		Parent: nil,
		Key:    key,
		Values: []any{value},
		Left:   nil,
		Right:  nil,
	})
}

func (t *BinaryNode) Push(tree *BinaryTree, node *BinaryNode) {
	switch tree.Compare(t.Key, node.Key) {
	case 0: // Equal
		t.Values = append(t.Values, node.Values...)

	case 1: // More
		if t.Right != nil {
			if t.Left == nil && !t.Right.Empty() {

			}
		}

	case -1: // Less

	}
}

func (t *BinaryNode) Empty() bool {
	return t.Right != nil || t.Left != nil
}

func Swap(l *BinaryNode, r *BinaryNode) {
	l.Parent, r.Parent = r.Parent, l.Parent
	l.Key, r.Key = r.Key, l.Key
	l.Values, r.Values = r.Values, l.Values
	l.Left, r.Left = r.Left, l.Left
}
