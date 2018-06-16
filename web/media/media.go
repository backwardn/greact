package media

import (
	"github.com/gopherjs/gopherjs/js"
)

type MediaQueryList interface {
	AddListener(func(*js.Object))
	RemoveListener(func(*js.Object))
	Matches() bool
}

type Event struct {
	*js.Object
	Matches bool `js:"matches"`
}

type Query struct {
	Query           string
	IsUnconditional bool
	mql             MediaQueryList
	handlers        []*Handler
}

func NewMediaQuery(mql MediaQueryList, query string, isUnconditional bool) *Query {
	m := &Query{
		Query:           query,
		IsUnconditional: isUnconditional,
	}
	m.mql = mql
	m.mql.AddListener(m.listen)
	return m
}

func (m *Query) clear() {
	for _, v := range m.handlers {
		v.destroy()
	}
	m.mql.RemoveListener(m.listen)
	m.handlers = nil
}

func (m *Query) listen(o *js.Object) {
	e := &Event{Object: o}
	var on bool
	if e.Matches || m.IsUnconditional {
		on = true
	}
	for _, v := range m.handlers {
		if on {
			v.on()
		} else {
			v.off()
		}
	}
}

func (m *Query) AddHandler(h *Handler) {
	m.handlers = append(m.handlers, h)
	if m.mql.Matches() || m.IsUnconditional {
		h.on()
	}
}

type Options struct {
	Match      func()
	UnMatch    func()
	Setup      func()
	Destroy    func()
	DeferSetup bool
}

type Handler struct {
	options     *Options
	initialized bool
}

func NewQueryHandler(opts *Options) *Handler {
	q := &Handler{options: opts}
	if !opts.DeferSetup {
		q.setup()
	}
	return q
}

func (q *Handler) setup() {
	if q.options.Setup != nil {
		q.options.Setup()
	}
	q.initialized = true
}

func (q *Handler) on() {
	if !q.initialized {
		q.setup()
	}
	if q.options.Match != nil {
		q.options.Match()
	}
}

func (q *Handler) off() {
	if q.options.UnMatch != nil {
		q.options.UnMatch()
	}
}

func (q *Handler) destroy() {
	if q.options.Destroy != nil {
		q.options.Destroy()
	} else {
		q.off()
	}
}

type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}

var ErrNotSupported = &Error{msg: "matchMedia not present, legacy browsers require a polyfill"}

type Dispatch struct {
	BrowserIsIncapable bool
	queries            map[string]*Query
}

func NewDispatch(isIncapable bool) *Dispatch {
	// m := js.Global.Get("matchMedia")
	// if m == nil {
	// 	panic(ErrNotSupported)
	// }
	// s := js.Global.Call("matchMedia", "only all").Get("matches").Bool()
	return &Dispatch{BrowserIsIncapable: isIncapable}
}

func (d *Dispatch) Register(mql MediaQueryList, query string, shoudDegrade bool, opts ...*Options) {
	isUnconditional := shoudDegrade && d.BrowserIsIncapable
	q, ok := d.queries[query]
	if !ok {
		q = NewMediaQuery(mql, query, isUnconditional)
		d.queries[query] = q
	}
	for _, v := range opts {
		q.AddHandler(NewQueryHandler(v))
	}
}

func (d *Dispatch) UnRegister(query string) {
	if q, ok := d.queries[query]; ok {
		q.clear()
	}
}
