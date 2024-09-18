//go:build darwin

/* -----------------------------------------------------------------
 *				L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 LordOfScripts
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * MacOS-specific Firefox
 *-----------------------------------------------------------------*/
package firefox

import (
	"path/filepath"

	cmn "github.com/lordofscripts/wipechromium"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	VARIANT string = "Firefox"
)

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// NOTE: 'profileDir' here is not its name but its actual sub-path directory!
func GetFirefoxDirs(profileDir string) (error, string, string) {
	_, dataDir := GetDataDir(profileDir)
	_, cacheDir := GetCacheDir(profileDir)
	return nil, dataDir, cacheDir
}

func GetRootDataDir() (error, string) {
	subPath := filepath.Join("Library", "Caches", "Firefox")
	return nil, cmn.AtHome(subPath)
}

// Profile-specific Cache directory
// *MacOS: ~/Users/USERNAME/Library/Caches/Firefox/Profiles/PROFILE/{cache|cache2}/
func GetCacheDir(profileDir string) (error, string) {
	subPath := filepath.Join("Library", "Caches", "Firefox", "Profiles", profileDir, "cache2")
	return nil, cmn.AtHome(subPath)
}

// Profile-specific Profile directory
// *MacOS: ~/Users/USERNAME/Library/Caches/Firefox/Profiles/PROFILE/
func GetDataDir(profileDir string) (error, string) {
	subPath := filepath.Join("Library", "Caches", "Firefox", "Profiles", profileDir)
	return nil, cmn.AtHome(subPath)
}
