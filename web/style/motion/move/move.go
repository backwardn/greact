package move

import (
	"github.com/gernest/vected/lib/gs"
	"github.com/gernest/vected/web/style/core/themes"
	"github.com/gernest/vected/web/style/mixins"
)

func Motion(klass, keyframe string) gs.CSSRule {
	return gs.CSS(
		mixins.MakeMotion(klass, keyframe, themes.Default.AnimationDurationBase),
		gs.S(klass+"-enter",
			gs.S(""+klass+"-appear",
				gs.P("opacity", "0"),
				gs.P("animation-timing-function", themes.Default.EaseOutCirc),
			),
		),
		gs.S(klass+"-leave",
			gs.P("animation-timing-function", themes.Default.EaseInCirc),
		),
	)
}

//keyframes
const (
	Down  = "moveDown"
	Up    = "moveUp"
	Left  = "moveLeft"
	Right = "moveRight"
)

func KeyFrames() gs.CSSRule {
	return gs.CSS(
		gs.KeyFrame(Down+"In",
			gs.Cond("0%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateY(100%)"),
				gs.P("opacity", "0"),
			),
			gs.Cond("100%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateY(0%)"),
				gs.P("opacity", "1"),
			),
		),
		gs.KeyFrame(Down+"Out",
			gs.Cond("0%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateY(100%)"),
				gs.P("opacity", "1"),
			),
			gs.Cond("100%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateY(0%)"),
				gs.P("opacity", "0"),
			),
		),
		gs.KeyFrame(Left+"In",
			gs.Cond("0%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateX(-100%)"),
				gs.P("opacity", "0"),
			),
			gs.Cond("100%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateX(0%)"),
				gs.P("opacity", "1"),
			),
		),
		gs.KeyFrame(Left+"Out",
			gs.Cond("0%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateX(0%)"),
				gs.P("opacity", "1"),
			),
			gs.Cond("100%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateX(-100%)"),
				gs.P("opacity", "0"),
			),
		),
		gs.KeyFrame(Right+"In",
			gs.Cond("0%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateX(100%)"),
				gs.P("opacity", "0"),
			),
			gs.Cond("100%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateX(0%)"),
				gs.P("opacity", "1"),
			),
		),
		gs.KeyFrame(Right+"Out",
			gs.Cond("0%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateX(0%)"),
				gs.P("opacity", "1"),
			),
			gs.Cond("100%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateX(100%)"),
				gs.P("opacity", "0"),
			),
		),
		gs.KeyFrame(Up+"In",
			gs.Cond("0%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateY(-100%)"),
				gs.P("opacity", "0"),
			),
			gs.Cond("100%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateY(0%)"),
				gs.P("opacity", "1"),
			),
		),
		gs.KeyFrame(Up+"Out",
			gs.Cond("0%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateY(0%)"),
				gs.P("opacity", "1"),
			),
			gs.Cond("100%",
				gs.P("transform-origin", "0 0"),
				gs.P("transform", "translateY(-100%)"),
				gs.P("opacity", "0"),
			),
		),
	)
}
