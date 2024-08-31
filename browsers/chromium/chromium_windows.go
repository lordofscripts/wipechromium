//go:build windows

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Windows-specific Chromium
 *-----------------------------------------------------------------*/
package chromium

import (
	"os"
	"path/filepath"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

func GetChromiumDirs() (string, string) {
	return GetDataDir(), GetCacheDir()
}

// Profile-specific Cache directory
// *Unix/Linux: ~/.cache/chromium/Default
// *MacOS: ~/Library/Caches/Google/Chromium/Default
// *Windows:
func GetCacheDir() string { // TODO: (Windows) needs to be verified!
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err.Error())
	}
	ChromiumCachesDir := filepath.Join(cacheDir, "Chromium", "User Data", "Default", "Cache")

	return ChromiumCachesDir
}

// Profile-specific Profile directory
// *Unix/Linux: ~/.config/chromium
// *MacOS: ~/Library/Application Support/Chromium/Default
// *Windows: %LOCALAPPDATA%\Google\Chrome\User Data\Default
func GetDataDir() string { // TODO: (Windows) needs to be verified!
	cacheDir := GetCacheDir()
	ChromiumProfilesDir := filepath.Join(cacheDir, "Chromium", "User Data", "Default")
	return ChromiumProfilesDir
}
