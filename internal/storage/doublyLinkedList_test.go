package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_doubly_linked_list(t *testing.T) {
	for scenario, fn := range map[string]func(
		t *testing.T, dll dll, nd []*node,
	){
		"node":        testNode,
		"pop":         testPopNode,
		"shiftNode":   testShiftNode,
		"unshift":     testUnshift,
		"unshiftNode": testUnshiftNode,
		"removeNode":  testRemoveNode,
	} {
		t.Run(scenario, func(t *testing.T) {
			var nodes []*node
			node, err := NewNode("key0", "value0")
			require.NoError(t, err)
			nodes = append(nodes, node)
			node, err = NewNode("key1", "value1")
			require.NoError(t, err)
			nodes = append(nodes, node)
			dll := NewDll(3)
			fn(t, dll, nodes)
		})
	}
}

func testNode(t *testing.T, dll dll, nd []*node) {
	require.Equal(t, "key0", nd[0].key)
	require.Equal(t, "value0", nd[0].val)
	require.Nil(t, nd[0].next)
	require.Nil(t, nd[0].prev)
}

func testRemoveNode(t *testing.T, dll dll, nd []*node) {
	// case node is the tail
	dll.unshiftNode(nd[0])
	dll.unshiftNode(nd[1])

	rmNode := dll.removeNode(nd[1])
	require.Equal(t, rmNode.val, "value1")
	require.Equal(t, dll.length, 1)

	// case node is the head
	rmNode = dll.removeNode(nd[0])
	require.Equal(t, rmNode.val, "value0")
	require.Equal(t, dll.length, 0)

	// else
	dll.unshiftNode(nd[0])
	dll.unshiftNode(nd[1])
	node3, _ := NewNode("key2", "value2")
	dll.unshiftNode(node3)
	rmNode = dll.removeNode(nd[1])
	require.Equal(t, rmNode.val, "value1")
	require.Equal(t, dll.length, 2)
	require.Equal(t, dll.head.next.val, "value0")
	require.Nil(t, dll.head.prev)
	require.Equal(t, dll.tail.prev.val, "value2")
	require.Nil(t, dll.tail.next)
}

func testUnshiftNode(t *testing.T, dll dll, nd []*node) {
	// case dll.length == 0
	insertedNode, outboundedNode := dll.unshiftNode(nd[0])
	require.Equal(t, insertedNode.val, "value0")
	require.Nil(t, outboundedNode)

	// case dll.length > 0 and not over dll.length
	insertedNode, outboundedNode = dll.unshiftNode(nd[1])
	require.Equal(t, insertedNode.val, "value1")
	require.Nil(t, outboundedNode)

	dll.maxLgt = 2
	node3, _ := NewNode("key2", "value2")
	// case dll.length > 0 and over dll.length
	insertedNode, outboundedNode = dll.unshiftNode(node3)
	require.Equal(t, insertedNode.val, "value2")
	require.Equal(t, outboundedNode.val, "value0")
}

func testUnshift(t *testing.T, dll dll, nd []*node) {
	// case dll.length == 0
	insertedNode, outboundedNode, err := dll.unshift("key0", "value0")
	require.Equal(t, insertedNode.val, "value0")
	require.Nil(t, outboundedNode)
	require.Nil(t, err)

	// case dll.length > 0 and not over dll.length
	insertedNode, outboundedNode, err = dll.unshift("key1", "value1")
	require.Equal(t, insertedNode.val, "value1")
	require.Nil(t, outboundedNode)
	require.Nil(t, err)

	dll.maxLgt = 2
	// case dll.length > 0 and over dll.length
	insertedNode, outboundedNode, err = dll.unshift("key2", "value2")
	require.Equal(t, insertedNode.val, "value2")
	require.Equal(t, outboundedNode.val, "value0")
	require.Nil(t, err)
}

func testPopNode(t *testing.T, dll dll, nd []*node) {
	// case dll empty
	node := dll.popNode()
	require.Nil(t, node)

	dll.unshiftNode(nd[0])
	dll.unshiftNode(nd[1])

	// case dll.length > 1
	node = dll.popNode()
	require.Equal(t, "value0", node.val)
	require.Equal(t, dll.length, 1)
	require.Equal(t, dll.head.val, "value1")
	require.Equal(t, dll.tail.val, "value1")

	// case dll.length == 1
	node = dll.popNode()
	require.Equal(t, "value1", node.val)
	require.Equal(t, dll.length, 0)
	require.Nil(t, dll.head)
	require.Nil(t, dll.tail)

}

func testShiftNode(t *testing.T, dll dll, nd []*node) {
	// case dll empty
	node := dll.shiftNode()
	require.Nil(t, node)

	dll.unshiftNode(nd[0])
	dll.unshiftNode(nd[1])

	// case dll.length > 1
	node = dll.shiftNode()
	require.Equal(t, "value1", node.val)
	require.Equal(t, dll.length, 1)
	require.Equal(t, dll.head.val, "value0")
	require.Equal(t, dll.tail.val, "value0")

	// case dll.length == 1
	node = dll.popNode()
	require.Equal(t, "value0", node.val)
	require.Equal(t, dll.length, 0)
	require.Nil(t, dll.head)
	require.Nil(t, dll.tail)
}
