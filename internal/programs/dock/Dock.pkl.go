// Code generated from Pkl module `macfigure.dock`. DO NOT EDIT.
package dock

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Dock struct {
	Apps []any `pkl:"apps"`

	Folders []string `pkl:"folders"`

	AppswitcherAllDisplays *bool `pkl:"appswitcher-all-displays"`

	Autohide *bool `pkl:"autohide"`

	AutohideDelay *float64 `pkl:"autohide-delay"`

	AutohideTimeModifier *float64 `pkl:"autohide-time-modifier"`

	DashboardInOverlay *bool `pkl:"dashboard-in-overlay"`

	EnableSpringLoadActionsOnAllItems *bool `pkl:"enable-spring-load-actions-on-all-items"`

	ExposeAnimationDuration *float64 `pkl:"expose-animation-duration"`

	ExposeGroupApps *bool `pkl:"expose-group-apps"`

	Launchanim *bool `pkl:"launchanim"`

	Mineffect *string `pkl:"mineffect"`

	MinimizeToApplication *bool `pkl:"minimize-to-application"`

	MouseOverHiliteStack *bool `pkl:"mouse-over-hilite-stack"`

	MruSpaces *bool `pkl:"mru-spaces"`

	Orientation *string `pkl:"orientation"`

	ScrollToOpen *bool `pkl:"scroll-to-open"`

	ShowAppExposeGestureEnabled *bool `pkl:"showAppExposeGestureEnabled"`

	ShowDesktopGestureEnabled *bool `pkl:"showDesktopGestureEnabled"`

	ShowLaunchpadGestureEnabled *bool `pkl:"showLaunchpadGestureEnabled"`

	ShowMissionControlGestureEnabled *bool `pkl:"showMissionControlGestureEnabled"`

	ShowProcessIndicators *bool `pkl:"show-process-indicators"`

	Showhidden *bool `pkl:"showhidden"`

	ShowRecents *bool `pkl:"show-recents"`

	SlowMotionAllowed *bool `pkl:"slow-motion-allowed"`

	StaticOnly *bool `pkl:"static-only"`

	Tilesize *int `pkl:"tilesize"`

	Magnification *bool `pkl:"magnification"`

	Largesize *int `pkl:"largesize"`

	WvousTlCorner *int `pkl:"wvous-tl-corner"`

	WvousBlCorner *int `pkl:"wvous-bl-corner"`

	WvousTrCorner *int `pkl:"wvous-tr-corner"`

	WvousBrCorner *int `pkl:"wvous-br-corner"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Dock
func LoadFromPath(ctx context.Context, path string) (ret Dock, err error) {
	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions)
	if err != nil {
		return ret, err
	}
	defer func() {
		cerr := evaluator.Close()
		if err == nil {
			err = cerr
		}
	}()
	ret, err = Load(ctx, evaluator, pkl.FileSource(path))
	return ret, err
}

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Dock
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Dock, error) {
	var ret Dock
	err := evaluator.EvaluateModule(ctx, source, &ret)
	return ret, err
}
