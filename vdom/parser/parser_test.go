package parser

import (
	"testing"

	"github.com/gernest/vected/vdom"
)

func TestClear(t *testing.T) {
	t.Run("should return  element", func(ts *testing.T) {
		e := `<div></div>`
		n, err := ParseString(e)
		if err != nil {
			ts.Fatal(err)
		}
		if n.Data != "div" {
			t.Errorf("expected div got %s", n.Data)
		}
	})
	t.Run("should return  container element", func(ts *testing.T) {
		e := `
		<div>
		</div>
		<div>
		</div>
		`
		n, err := ParseString(e)
		if err != nil {
			ts.Fatal(err)
		}
		if n.Data != vdom.ContainerNode {
			t.Errorf("expected %s got %s", vdom.ContainerNode, n.Data)
		}
	})
}

func TestGenerate(t *testing.T) {
	n, err := ParseString(`<div className={props.classNames}></div>`)
	if err != nil {
		t.Fatal(err)
	}
	v, err := GenerateRenderMethod(n, &Context{
		Package:    "test",
		Recv:       "h",
		StructName: "Hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	// ioutil.WriteFile("test/test.gen.go", v, 0600)
	s := string(v)
	if s != expected1 {
		t.Errorf("got wrong generated output")
	}
}

const expected1 = `// Code generated by vected DO NOT EDIT.
package test

import (
	"context"
	"github.com/gernest/vected/props"
	"github.com/gernest/vected/state"
	"github.com/gernest/vected/vdom"
)

// Render implements vected.Renderer interface.
func (h Hello) Render(ctx context.Context, props props.Props, state state.State) *vdom.Node {
	return &vdom.Node{
		Type: vdom.ElementNode,
		Data: "div",
		Attr: []vdom.Attribute{
			{Key: "classname", Val: props["classNames"]},
		},
	}
}
`