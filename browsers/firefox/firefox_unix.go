//go:build linux || aix || freebsd || netbsd || openbsd || solaris

/* -----------------------------------------------------------------
 *				L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 LordOfScripts
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Linux/Unix-specific Chromium
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
	VARIANT string = "Firefox-ESR"
)

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// Returns the Data & Cache directories which are NOT profile-specific.
// The profile name still has to be added for each specific profile.
// NOTE: In Linux/Un*x we are already working within a user account dir
// NOTE: 'profileDir' here is not its name but its actual sub-path directory!
func GetFirefoxDirs(profileDir string) (error, string, string) {
	_, dataDir := GetDataDir(profileDir)
	_, cacheDir := GetCacheDir(profileDir)
	return nil, dataDir, cacheDir
}

func GetRootDataDir() (error, string) {
	return nil, cmn.AtHome(filepath.Join(".mozilla", "firefox"))
}

// Profile-specific Cache directory
// *Unix/Linux: ~/.cache/mozilla/firefox/
func GetCacheDir(profileDir string) (error, string) {
	return nil, cmn.AtHome(filepath.Join(".cache", "mozilla", "firefox", profileDir))
}

// Profiles-specific Profiles directory
// *Unix/Linux: ~/.mozilla/
func GetDataDir(profileDir string) (error, string) {
	return nil, cmn.AtHome(filepath.Join(".mozilla", "firefox", profileDir))
}
