// Package vected is a component based frontend framework for golang. This
// framework delivers high performance and sleek ui's, that works both on the
// serverside and the frontend.
//
// Go templates are used as the main templating system. The framework is
// inspired by react, especially preact which I used to learn more about how
// react works.
//
// Also, this borrows from vue js. The templates are just normal go templates so
// no need to learn a different syntax.
//
// The user intrface styles is ant design
// see https://github.com/ant-design/ant-design to learn more about ant design.
package vected

import (
	"container/list"
	"context"
	"sync"

	"github.com/gernest/vected/vdom/value"

	"github.com/gernest/vected/prop"
	"github.com/gernest/vected/state"
	"github.com/gernest/vected/vdom"
	"github.com/gernest/vected/vdom/dom"
)

// RenderMode is a flag determining how a component is rendered.
type RenderMode uint

//supported render mode
const (
	No RenderMode = iota
	Force
	Sync
	Async
)

var queue = NeqQueueRenderer()
var mounts = list.New()

// Component is an interface which defines a unit of user interface.
type Component interface {
	Template() string
	Render(context.Context, prop.Props, state.State) *vdom.Node
	core() *Core
}

// Core is th base struct that every struct that wants to implement Component
// interface must embed.
//
// This is used to make Props available to the component.
type Core struct {
	props           prop.Props
	state           state.State
	prevProps       prop.Props
	prevState       state.State
	disable         bool
	renderCallbacks []func()
	context         context.Context
	prevContext     context.Context
	component       Component
	base            dom.Element
	nextBase        dom.Element
	dirty           bool
	key             prop.NullString

	// This is a callback used to receive instance of Component or the Dom element.
	// after they have been mounted.
	ref func(interface{})

	// priority this is a number indicating how important this component is in the
	// re rendering queue. The higher the number the more urgent re renders.
	priority int
}

func (c *Core) core() *Core { return c }

// SetState updates component state and schedule re rendering.
func (c *Core) SetState(newState state.State, callback ...func()) {
	prev := c.prevState
	c.prevState = newState
	c.state = state.Merge(prev, newState)
	if len(callback) > 0 {
		c.renderCallbacks = append(c.renderCallbacks, callback...)
	}
	//TODO enqueue this for re rendering.
}

// InitState is an interface for exposing initial state.
// Component should implement this interface if they want to set initial state
// when the component is first created before being rendered.
type InitState interface {
	InitState() state.State
}

// InitProps is an interface for exposing default props. This will be merged
// with other props before being sent to render.
type InitProps interface {
	InitProps() prop.Props
}

// WillMount is an interface defining a callback which is invoked before the
// component is mounted on the dom.
type WillMount interface {
	ComponentWillMount()
}

// DidMount is an interface defining a callback that is invoked after the
// component has been mounted to the dom.
type DidMount interface {
	ComponentDidMount()
}

// WillUnmount is an interface defining a callback that is invoked prior to
// removal of the rendered component from the dom.
type WillUnmount interface {
	ComponentWillUnmount()
}

// WillReceiveProps is an interface defining a callback that will be called with
// the new props before they are accepted and passed to be rendered.
type WillReceiveProps interface {
	ComponentWillReceiveProps(context.Context, prop.Props)
}

// ShouldUpdate is an interface defining callback that is called before render
// determine if re render is necessary.
type ShouldUpdate interface {
	// If this returns false then re rendering for the component is skipped.
	ShouldComponentUpdate(context.Context, prop.Props, state.State) bool
}

// WillUpdate is an interface defining a callback that is called before rendering
type WillUpdate interface {
	// If returned props are not nil, then it will be merged with nextprops then
	// passed to render for rendering.
	ComponentWillUpdate(context.Context, prop.Props, state.State) prop.Props
}

// DidUpdate defines a callback that is invoked after rendering.
type DidUpdate interface {
	ComponentDidUpdate()
}

// DerivedState is an interface which can be used to derive state from props.
type DerivedState interface {
	DeriveState(prop.Props, state.State) state.State
}

// SetProps sets cmp props and possibly re renders
func SetProps(ctx context.Context, cmp Component, props prop.Props, state state.State, mode RenderMode, mountAll bool) {
	core := cmp.core()
	if core.disable {
		return
	}
	ref := props["ref"]
	if fn, ok := ref.(func(interface{})); ok {
		core.ref = fn
	}
	core.key = props.String("key")
	delete(props, "key")
	delete(props, "ref")
	_, ok := cmp.(DerivedState)
	if !ok {
		if core.base == nil || mountAll {
			if m, ok := cmp.(WillMount); ok {
				m.ComponentWillMount()
			}
		} else if m, ok := cmp.(WillReceiveProps); ok {
			m.ComponentWillReceiveProps(ctx, props)
		}
	}
	if core.prevProps == nil {
		core.prevProps = core.props
	}
	core.props = props
	core.disable = false
	if mode != No {
		if mode == Sync {
			renderComponent(cmp, Sync, mountAll)
		} else {
			enqueueRender(cmp)
		}
	}
	if core.ref != nil {
		core.ref(cmp)
	}
}

func renderComponent(cmp Component, mode RenderMode, mountAll bool, child ...bool) {

}

type QueuedRender struct {
	components *list.List
	mu         sync.RWMutex
	closed     bool
}

func (q *QueuedRender) Push(v Component) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.components.PushBack(v)
}

// Pop returns the last added component and removes it from the queue.
func (q *QueuedRender) Pop() Component {
	e := q.pop()
	if e != nil {
		return e.Value.(Component)
	}
	return nil
}

func (q *QueuedRender) pop() *list.Element {
	e := q.last()
	q.mu.Lock()
	if e != nil {
		q.components.Remove(e)
	}
	q.mu.Unlock()
	return e
}

func (q *QueuedRender) last() *list.Element {
	q.mu.RLock()
	e := q.components.Back()
	q.mu.RUnlock()
	return e
}

// Last returns the last added component to the queue.
func (q *QueuedRender) Last() Component {
	e := q.last()
	if e != nil {
		return e.Value.(Component)
	}
	return nil
}

func NeqQueueRenderer() *QueuedRender {
	return &QueuedRender{
		components: list.New(),
	}
}

func enqueueRender(cmp Component) {
	if cmp.core().dirty {
		queue.Push(cmp)
		queue.Rerender()
	}
}

// Rerender re renders all enqueued dirty components.
func (q *QueuedRender) Rerender() {
	go q.rerender()
}

func (q *QueuedRender) rerender() {
	for cmp := q.Pop(); cmp != nil; cmp = q.Pop() {
		if cmp.core().dirty {
			renderComponent(cmp, 0, false)
		}
	}
}

func flushMounts() {
	for c := mounts.Back(); c != nil; c = mounts.Back() {
		if cmp, ok := c.Value.(Component); ok {
			if m, ok := cmp.(DidMount); ok {
				m.ComponentDidMount()
			}
		}
		mounts.Remove(c)
	}
}

func recollectNodeTree(node dom.Element, unmountOnly bool) {
	cmp := findComponent(node)
	if cmp != nil {
		unmountComponent(cmp)
	} else {
		if !unmountOnly {
			dom.RemoveNode(node)
		}
		removeChildren(node)
	}
}

// findComponent returns the component that rendered the node element. This
// returns nil if the node wasn't a component.
func findComponent(node dom.Element) Component {
	return nil
}

func unmountComponent(cmp Component) {
	core := cmp.core()
	core.disable = true
	base := core.base
	if wm, ok := cmp.(WillUnmount); ok {
		wm.ComponentWillUnmount()
	}
	core.base = nil
	if core.component != nil {
		unmountComponent(core.component)
	} else if base != nil {
		core.nextBase = base
		dom.RemoveNode(base)
		removeChildren(base)
	}
}

func removeChildren(node dom.Element) {
	node = node.Get("lastChild")
	for {
		if !dom.Valid(node) {
			return
		}
		next := node.Get("previousSibling")
		recollectNodeTree(node, true)
		node = next
	}
}

// UndefinedFunc is a function  that returns a javascript undefined value.
type UndefinedFunc func() value.Value

// Undefined is a work around to allow the library to work with/without wasm
// support.
//
// TODO: find a better way to handle this.
var Undefined UndefinedFunc

// Callback this is supposed to be defined by the package consumers.
var Callback dom.CallbackGenerator

func diffAttributes(node dom.Element, attrs, old []vdom.Attribute, isSvgMode bool) {
	a := mapAtts(attrs)
	b := mapAtts(old)
	for k, v := range b {
		if _, ok := a[k]; !ok {
			dom.SetAccessor(Callback, node, k, v, Undefined(), isSvgMode)
		}
	}
	for k := range a {
		switch k {
		case "children", "innerHTML":
			continue
		default:
			dom.SetAccessor(Callback, node, k, b[k], a[k], isSvgMode)
		}
	}
}

func mapAtts(attrs []vdom.Attribute) map[string]vdom.Attribute {
	m := make(map[string]vdom.Attribute)
	for _, v := range attrs {
		m[v.Key] = v
	}
	return m
}
