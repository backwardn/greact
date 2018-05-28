package integration

import (
	"github.com/gernest/mad"
	"github.com/gernest/mad/ws"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// Integration wraps mad.Component into a vecty component. This will render the
// mad.Component and execute the tests after being mounted.
type Integration struct {
	vecty.Core
	UUID      string
	Pkg       string
	Component *mad.Component
}

// Mount runs the integration tests and reports results via websocket.
func (c *Integration) Mount() {
	go func() {
		w, err := ws.New(c.UUID)
		if err != nil {
			panic(err)
		}
		v := mad.Exec(c.Component.Cases)
		err = w.Report(v, c.Pkg, c.UUID)
		if err != nil {
			println(err)
		}
	}()
}

// Render implements vecty.Component interface. This works under the assumption
// the Component field call returns a vecty.ComponentOrHTML
func (c *Integration) Render() vecty.ComponentOrHTML {
	if c.Component.IsBody {
		return c.Component.Component().(vecty.ComponentOrHTML)
	}
	return elem.Body(
		c.Component.Component().(vecty.ComponentOrHTML),
	)
}