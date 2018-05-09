package prom

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gopherjs/gopherjs/js"
)

var (
	_ Test   = (*T)(nil)
	_ Test   = (*Suite)(nil)
	_ Test   = (*ExecCommand)(nil)
	_ Test   = (List)(nil)
	_ Result = (*baseResult)(nil)
	_ Result = (*RSWithNode)(nil)
)

type Test interface {
	run()
}

func Describe(desc string, ctx ...Test) Test {
	return &Suite{Desc: desc, Cases: List(ctx)}
}

type List []Test

func (ls List) run() {}

type Suite struct {
	Desc     string
	Cases    List
	ResultFN func() Result
}

func defaultResultFn() Result {
	return &baseResult{}
}

func (*Suite) run() {}

func It(desc string, fn func(Result)) Test {
	return &ExecCommand{Desc: desc, Func: fn}
}

type Result interface {
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	FatalF(string, ...interface{})
	Errors() []error
}

type ExecCommand struct {
	Desc string
	Func func(Result)
}

func (*ExecCommand) run() {}

type ResultInfo struct {
	Case         string   `json:"case"`
	Failed       bool     `json:"failed"`
	FailMessages []string `json:"fail_messages"`
}

type Error struct {
	Message error
}

func (e *Error) Error() string {
	if e.Message != nil {
		return e.Error()
	}
	return ""
}

type ResultCtx struct {
	Parent   *ResultCtx    `json:"-"`
	Desc     string        `json:"description"`
	Children []*ResultCtx  `json:"children"`
	Results  []*ResultInfo `json:"results"`
}

func (r ResultCtx) ToJson() string {
	v, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(v)
}

func Exec(ctx ...*T) *ResultCtx {
	rs := &ResultCtx{}
	for _, v := range ctx {
		rs.Children = append(rs.Children, v.exec())
	}
	return rs
}

func ExecSuite(s *Suite) *ResultCtx {
	rs := &ResultCtx{
		Desc: s.Desc,
	}
	for _, v := range s.Cases {
		switch e := v.(type) {
		case *Suite:
			ch := ExecSuite(e)
			ch.Parent = rs
			rs.Children = append(rs.Children, ch)
		case *ExecCommand:
			fn := s.ResultFN
			if fn == nil {
				fn = defaultResultFn
			}
			rv, ok := wrapPanic(e, fn)
			rs.Results = append(rs.Results, rv)
			if ok {
				// call to Fatal or Fatalf halts the whole suite.
				return rs
			}
		}
	}
	return rs
}

func wrapPanic(e *ExecCommand, fn func() Result) (rs *ResultInfo, panicked bool) {
	defer func() {
		if err := recover(); err != nil {
			rs.Failed = true
			rs.FailMessages = append(rs.FailMessages, fmt.Sprint(err))
			panicked = true
		}
	}()
	rs = execute(e, fn)
	return
}

type baseResult struct {
	err []error
}

func (b *baseResult) Error(v ...interface{}) {
	b.err = append(b.err, errors.New(fmt.Sprint(v...)))
}
func (b *baseResult) Fatal(v ...interface{}) {
	panic(&Error{Message: errors.New(fmt.Sprint(v...))})
}

func (b *baseResult) Errorf(s string, v ...interface{}) {
	b.err = append(b.err, fmt.Errorf(s, v...))
}

func (b *baseResult) FatalF(s string, v ...interface{}) {
	panic(&Error{Message: fmt.Errorf(s, v...)})
}

func (b *baseResult) Errors() []error {
	return b.err
}

// execute calls the function e.Func and register results.
func execute(e *ExecCommand, fn func() Result) (rs *ResultInfo) {
	r := fn()
	if e.Func != nil {
		e.Func(r)
	}
	rs = &ResultInfo{Case: e.Desc}
	errs := r.Errors()
	if errs != nil {
		rs.Failed = true
		for _, v := range errs {
			rs.FailMessages = append(rs.FailMessages, v.Error())
		}
	}
	return rs
}

type T struct {
	before func()
	after  func(*ResultCtx)
	suit   *Suite
	base   List
}

func NewTest(name string, before func(), after func(*ResultCtx)) *T {
	return &T{
		suit:   &Suite{Desc: name},
		before: before,
		after:  after,
	}
}

func (t *T) Before(fn ...func()) {
	if len(fn) > 0 {
		t.before = fn[0]
	}
}

func (t *T) After(fn ...func(*ResultCtx)) {
	if len(fn) > 0 {
		t.after = fn[0]
	}
}

func (t *T) run() {}

func (t *T) Describe(desc string, cases ...Test) {
	t.suit.Cases = append(t.suit.Cases, Describe(desc, cases...))
}

func (t *T) exec() *ResultCtx {
	if t.before != nil {
		t.before()
	}
	rs := ExecSuite(t.suit)
	if t.after != nil {
		t.after(rs)
	}
	return rs
}

func (t *T) Cases(tc ...Test) *T {
	t.suit.Cases = append(t.suit.Cases, tc...)
	return t
}

type Component struct {
	ID        string
	Component func() interface{}
	IsBody    bool
	Cases     List
}

func (c *Component) runIntegration() {}

type Integration interface {
	runIntegration()
}

// Node is an interface for retrieving a rendered Component node.
type Node interface {
	Node() *js.Object
}

type RSWithNode struct {
	baseResult
	Object *js.Object
}

func (rs *RSWithNode) Node() *js.Object {
	return rs.Object
}

// Render returns an integration test for non body Components. Use this to test
// Components that renders spans,div etc.
//
// NOTE: func()interface{} was supposed to be func()vecty.ComponentOrHTML this
// is a workaround, because importing vecty in this package will make it
// impossible to run the commandline tools since vecty only works with the
// browser.
func Render(desc string, c func() interface{}, cases ...Test) Integration {
	return &Component{
		ID: desc, Component: c, Cases: cases,
	}
}

// RenderBody is like Render but the Component is expected to be elem.Body
func RenderBody(desc string, c func() interface{}, cases ...Test) Integration {
	return &Component{
		ID: desc, Component: c, Cases: cases, IsBody: true,
	}
}