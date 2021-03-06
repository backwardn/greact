// package greact is a component based frontend framework for golang. This
// framework delivers high performance and responsive ui.
//
// This relies on the experimental wasm api to interact with dom. The project
// started as a port of preact to go, but has since evolved. It still borrows a
// similar API from react/preact.
package greact

import "github.com/gernest/greact/node"

type Core interface {
	core()
}

type Props = node.Props

type State = node.State
