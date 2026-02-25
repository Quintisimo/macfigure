// Code generated from Pkl module `macfigure.nsglobaldomain`. DO NOT EDIT.
package nsglobaldomain

import (
	"context"

	"github.com/apple/pkl-go/pkl"
)

type Nsglobaldomain struct {
	AppleShowAllFiles *bool `pkl:"AppleShowAllFiles"`

	AppleEnableMouseSwipeNavigateWithScrolls *bool `pkl:"AppleEnableMouseSwipeNavigateWithScrolls"`

	AppleEnableSwipeNavigateWithScrolls *bool `pkl:"AppleEnableSwipeNavigateWithScrolls"`

	AppleFontSmoothing *int `pkl:"AppleFontSmoothing"`

	AppleInterfaceStyle *string `pkl:"AppleInterfaceStyle"`

	AppleIconAppearanceTheme *string `pkl:"AppleIconAppearanceTheme"`

	AppleInterfaceStyleSwitchesAutomatically *bool `pkl:"AppleInterfaceStyleSwitchesAutomatically"`

	AppleKeyboardUIMode *int `pkl:"AppleKeyboardUIMode"`

	ApplePressAndHoldEnabled *bool `pkl:"ApplePressAndHoldEnabled"`

	AppleShowAllExtensions *bool `pkl:"AppleShowAllExtensions"`

	AppleShowScrollBars *string `pkl:"AppleShowScrollBars"`

	AppleScrollerPagingBehavior *bool `pkl:"AppleScrollerPagingBehavior"`

	AppleSpacesSwitchOnActivate *bool `pkl:"AppleSpacesSwitchOnActivate"`

	NSAutomaticCapitalizationEnabled *bool `pkl:"NSAutomaticCapitalizationEnabled"`

	NSAutomaticInlinePredictionEnabled *bool `pkl:"NSAutomaticInlinePredictionEnabled"`

	NSAutomaticDashSubstitutionEnabled *bool `pkl:"NSAutomaticDashSubstitutionEnabled"`

	NSAutomaticPeriodSubstitutionEnabled *bool `pkl:"NSAutomaticPeriodSubstitutionEnabled"`

	NSAutomaticQuoteSubstitutionEnabled *bool `pkl:"NSAutomaticQuoteSubstitutionEnabled"`

	NSAutomaticSpellingCorrectionEnabled *bool `pkl:"NSAutomaticSpellingCorrectionEnabled"`

	NSAutomaticWindowAnimationsEnabled *bool `pkl:"NSAutomaticWindowAnimationsEnabled"`

	NSDisableAutomaticTermination *bool `pkl:"NSDisableAutomaticTermination"`

	NSDocumentSaveNewDocumentsToCloud *bool `pkl:"NSDocumentSaveNewDocumentsToCloud"`

	AppleWindowTabbingMode *string `pkl:"AppleWindowTabbingMode"`

	NSNavPanelExpandedStateForSaveMode *bool `pkl:"NSNavPanelExpandedStateForSaveMode"`

	NSNavPanelExpandedStateForSaveMode2 *bool `pkl:"NSNavPanelExpandedStateForSaveMode2"`

	NSTableViewDefaultSizeMode *int `pkl:"NSTableViewDefaultSizeMode"`

	NSTextShowsControlCharacters *bool `pkl:"NSTextShowsControlCharacters"`

	NSUseAnimatedFocusRing *bool `pkl:"NSUseAnimatedFocusRing"`

	NSScrollAnimationEnabled *bool `pkl:"NSScrollAnimationEnabled"`

	NSWindowResizeTime *float64 `pkl:"NSWindowResizeTime"`

	NSWindowShouldDragOnGesture *bool `pkl:"NSWindowShouldDragOnGesture"`

	NSStatusItemSpacing *int `pkl:"NSStatusItemSpacing"`

	NSStatusItemSelectionPadding *int `pkl:"NSStatusItemSelectionPadding"`

	KeyRepeat *int `pkl:"KeyRepeat"`

	PMPrintingExpandedStateForPrint *bool `pkl:"PMPrintingExpandedStateForPrint"`

	PMPrintingExpandedStateForPrint2 *bool `pkl:"PMPrintingExpandedStateForPrint2"`

	ComAppleKeyboardFnState *bool `pkl:"com.apple.keyboard.fnState"`

	ComAppleMouseTapBehavior *int `pkl:"com.apple.mouse.tapBehavior"`

	ComAppleSoundBeepVolume *float64 `pkl:"com.apple.sound.beep.volume"`

	ComAppleSoundBeepFeedback *int `pkl:"com.apple.sound.beep.feedback"`

	ComAppleTrackpadEnableSecondaryClick *bool `pkl:"com.apple.trackpad.enableSecondaryClick"`

	ComAppleTrackpadTrackpadCornerClickBehavior *int `pkl:"com.apple.trackpad.trackpadCornerClickBehavior"`

	ComAppleTrackpadScaling *float64 `pkl:"com.apple.trackpad.scaling"`

	ComAppleTrackpadForceClick *bool `pkl:"com.apple.trackpad.forceClick"`

	ComAppleSpringingEnabled *bool `pkl:"com.apple.springing.enabled"`

	ComAppleSpringingDelay *float64 `pkl:"com.apple.springing.delay"`

	ComAppleSwipescrolldirection *bool `pkl:"com.apple.swipescrolldirection"`

	AppleMeasurementUnits *string `pkl:"AppleMeasurementUnits"`

	AppleMetricUnits *int `pkl:"AppleMetricUnits"`

	AppleTemperatureUnit *string `pkl:"AppleTemperatureUnit"`

	AppleICUForce24HourTime *bool `pkl:"AppleICUForce24HourTime"`

	HIHideMenuBar *bool `pkl:"_HIHideMenuBar"`
}

// LoadFromPath loads the pkl module at the given path and evaluates it into a Nsglobaldomain
func LoadFromPath(ctx context.Context, path string) (ret Nsglobaldomain, err error) {
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

// Load loads the pkl module at the given source and evaluates it with the given evaluator into a Nsglobaldomain
func Load(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (Nsglobaldomain, error) {
	var ret Nsglobaldomain
	err := evaluator.EvaluateModule(ctx, source, &ret)
	return ret, err
}
