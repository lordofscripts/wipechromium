//go:build aix darwin dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris
/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Linux/Unix-specific Chromium
 *-----------------------------------------------------------------*/
package chromium

import (
	"path/filepath"

	cmn "lordofscripts/wipechromium"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	cCHROMIUM string = "chromium"
	cCHROME_CACHES string =".cache"
	cCHROME_PROFILES string =".config"
)

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// Returns the Data & Cache directories which are NOT profile-specific.
// The profile name still has to be added for each specific profile.
func GetChromiumDirs() (string, string) {
	return GetDataDir(), GetCacheDir()
}

// Profile-specific Cache directory
// *Unix/Linux: ~/.cache/chromium/
func GetCacheDir() string {
	ChromiumCachesDir := filepath.Join(cmn.AtHome(cCHROME_CACHES), cCHROMIUM)
	return ChromiumCachesDir
}

// Profiles-specific Profile directory
// *Unix/Linux: ~/.config/
func GetDataDir() string {
	ChromiumProfilesDir := filepath.Join(cmn.AtHome(cCHROME_PROFILES), cCHROMIUM)
	return ChromiumProfilesDir
}

