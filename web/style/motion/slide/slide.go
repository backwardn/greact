package slide

import (
	"github.com/gernest/vected/lib/gs"
	"github.com/gernest/vected/web/style/core/themes"
	"github.com/gernest/vected/web/style/mixins"
)

//keyframes names
const (
	Up    = "slideUp"
	Down  = "slideDown"
	Left  = "slideLeft"
	Right = "slideRight"
)

func KeyFrames() gs.CSSRule {
	return gs.CSS(
		gs.KeyFrame(Up+"In",
			gs.Cond("0%",
				gs.P("opacity", "0"),
				gs.P("transform-origin", "0% 0%"),
				gs.P("transform", "scaleY(.8)"),
			),
			gs.Cond("100%",
				gs.P("opacity", "1"),
				gs.P("transform-origin", "0% 0%"),
				gs.P("transform", "scaleY(1)"),
			),
		),
		gs.KeyFrame(Up+"Out",
			gs.Cond("0%",
				gs.P("opacity", "1"),
				gs.P("transform-origin", "0% 0%"),
				gs.P("transform", "scaleY(1)"),
			),
			gs.Cond("100%",
				gs.P("opacity", "0"),
				gs.P("transform-origin", "0% 0%"),
				gs.P("transform", "scaleY(.8)"),
			),
		),
		gs.KeyFrame(Down+"In",
			gs.Cond("0%",
				gs.P("opacity", "0"),
				gs.P("transform-origin", "100% 100%"),
				gs.P("transform", "scaleY(.8)"),
			),
			gs.Cond("100%",
				gs.P("opacity", "1"),
				gs.P("transform-origin", "100% 100%"),
				gs.P("transform", "scaleY(1)"),
			),
		),
		gs.KeyFrame(Down+"Out",
			gs.Cond("0%",
				gs.P("opacity", "1"),
				gs.P("transform-origin", "100% 100%"),
				gs.P("transform", "scaleY(1)"),
			),
			gs.Cond("100%",
				gs.P("opacity", "0"),
				gs.P("transform-origin", "100% 100%"),
				gs.P("transform", "scaleY(.8)"),
			),
		),
		gs.KeyFrame(Left+"In",
			gs.Cond("0%",
				gs.P("opacity", "0"),
				gs.P("transform-origin", "0% 0%"),
				gs.P("transform", "scaleX(.8)"),
			),
			gs.Cond("100%",
				gs.P("opacity", "1"),
				gs.P("transform-origin", "0% 0%"),
				gs.P("transform", "scaleX(1)"),
			),
		),
		gs.KeyFrame(Left+"Out",
			gs.Cond("0%",
				gs.P("opacity", "1"),
				gs.P("transform-origin", "0% 0%"),
				gs.P("transform", "scaleX(1)"),
			),
			gs.Cond("100%",
				gs.P("opacity", "0"),
				gs.P("transform-origin", "0% 0%"),
				gs.P("transform", "scaleX(.8)"),
			),
		),
		gs.KeyFrame(Left+"In",
			gs.Cond("0%",
				gs.P("opacity", "0"),
				gs.P("transform-origin", "100% 0%"),
				gs.P("transform", "scaleX(.8)"),
			),
			gs.Cond("100%",
				gs.P("opacity", "1"),
				gs.P("transform-origin", "100% 0%"),
				gs.P("transform", "scaleX(1)"),
			),
		),
		gs.KeyFrame(Left+"Out",
			gs.Cond("0%",
				gs.P("opacity", "1"),
				gs.P("transform-origin", "100% 0%"),
				gs.P("transform", "scaleX(1)"),
			),
			gs.Cond("100%",
				gs.P("opacity", "0"),
				gs.P("transform-origin", "100% 0%"),
				gs.P("transform", "scaleX(.8)"),
			),
		),
	)
}

func Motion(klass, keyframe string) gs.CSSRule {
	return gs.CSS(
		mixins.MakeMotion(klass, keyframe, themes.Default.AnimationDurationBase),
		gs.S(klass+"-enter",
			gs.S(""+klass+"-appear",
				gs.P("opacity", "0"),
				gs.P("animation-timing-function", themes.Default.EaseOutQuint),
			),
		),
		gs.S(klass+"-leave",
			gs.P("animation-timing-function", themes.Default.EaseInQuint),
		),
	)
}
