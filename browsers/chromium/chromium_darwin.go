//go:build darwin

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * MacOS-specific Chromium
 *-----------------------------------------------------------------*/
package chromium

import (
	"path/filepath"

	cmn "github.com/lordofscripts/wipechromium"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	cCHROMIUM        string = "Chromium"
	cCHROME_CACHES   string = "Library/Caches/Google"
	cCHROME_PROFILES string = "Library/Application Support"
)

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

func GetChromiumDirs() (string, string) {
	return GetDataDir(), GetCacheDir()
}

// Profile-specific Cache directory
// *Unix/Linux: ~/.cache/chromium/Default
// *MacOS: ~/Library/Caches/Google/Chromium/Default
func GetCacheDir() string {
	ChromiumCachesDir := filepath.Join(cmn.AtHome(cCHROME_CACHES), cCHROMIUM)
	return ChromiumCachesDir
}

// Profile-specific Profile directory
// *Unix/Linux: ~/.config/chromium
// *MacOS: ~/Library/Application Support/Chromium/Default
func GetDataDir() string {
	ChromiumProfilesDir := filepath.Join(cmn.AtHome(cCHROME_PROFILES), cCHROMIUM)
	return ChromiumProfilesDir
}
