package xmltree

import (
	"encoding/xml"
	"inttest-runtime/pkg/utils"
	"reflect"
	"strings"

	"github.com/samber/lo"
)

type nodeParser struct {
	Tag     xml.Name
	Content []byte       `xml:",innerxml"`
	Nodes   []nodeParser `xml:",any"`
}

func (p nodeParser) buildTree(parent *Node) *Node {
	root := &Node{
		Tag: XMLTag(lo.Ternary(
			p.Tag.Space == "",
			p.Tag.Local,
			p.Tag.Space+":"+p.Tag.Local,
		)),
		Content: utils.B2S(p.Content),
		Parent:  parent,
	}
	root.Children = lo.Map(p.Nodes, func(c nodeParser, _ int) *Node {
		return c.buildTree(root)
	})
	return root
}

type Node struct {
	Tag      XMLTag
	Content  string
	Children []*Node
	Parent   *Node
}

func (n *Node) Traversal(fn func(n *Node) error) error {
	if err := fn(n); err != nil {
		return err
	}
	for _, c := range n.Children {
		if err := c.Traversal(fn); err != nil {
			return err
		}
	}
	return nil
}

func (n *Node) Equal(other *Node) bool {
	// аче так можно было?
	return reflect.DeepEqual(n, other)
}

type XMLTag string

func (t XMLTag) String() string {
	return string(t)
}

func (t XMLTag) Namespace() string {
	tag := t.String()
	colonPos := strings.Index(tag, ":")
	if colonPos == -1 {
		return ""
	}
	return tag[:colonPos]
}

func (t XMLTag) Local() string {
	tag := t.String()
	colonPos := strings.Index(tag, ":")
	if colonPos == -1 {
		return tag
	}
	return tag[colonPos+1:]
}

func FromBytes(b []byte) (*Node, error) {
	var root nodeParser
	if err := xml.Unmarshal(b, &root); err != nil {
		return nil, err
	}

	return root.buildTree(nil), nil
}

func (n *Node) toParsingTree() nodeParser {
	return nodeParser{
		Tag: xml.Name{
			Space: n.Tag.Namespace(),
			Local: n.Tag.Local(),
		},
		Content: utils.S2B(n.Content),
		Nodes: lo.Map(n.Children, func(c *Node, _ int) nodeParser {
			return c.toParsingTree()
		}),
	}
}

func (n *Node) Marshal() ([]byte, error) {
	return xml.Marshal(n.toParsingTree())
}
