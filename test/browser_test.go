/* -----------------------------------------------------------------
 *				C o r a l y s   T e c h n o l o g i e s
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *						U n i t   T e s t
 *-----------------------------------------------------------------*/
package test

import (
	"lordofscripts/wipechromium/browsers"
	"testing"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

var (
	navigator browsers.Browser
)

/* ----------------------------------------------------------------
 *				U n i t  T e s t   F u n c t i o n s
 *-----------------------------------------------------------------*/
func Test_BrowserEnum(t *testing.T) {
	navigator = browsers.ChromiumBrowser

	if navigator.String() != "Chromium" {
		t.Errorf("unexpected String()")
	}
}

func Test_Invalid_BrowserEnum(t *testing.T) {
	navigator = 8

	if navigator.String() != "" {
		t.Errorf("Should be empty for unknown enum value")
	}
}
