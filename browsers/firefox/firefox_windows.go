//go:build windows

/* -----------------------------------------------------------------
 *				L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 LordOfScripts
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Windows-specific Chromium
 *-----------------------------------------------------------------*/
package firefox

import (
	"fmt"
	"os"
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

/*
	os.UserHomeDir() // C:\Users\YourUser
    os.UserCacheDir() // C:\Users\YourUser\AppData\Local
    os.UserConfigDir() // C:\Users\YourUser\AppData\Roaming
*/
// NOTE: In Windows we need the User Account folder AND the profile name
// NOTE: 'profileDir' here is not its name but its actual sub-path directory!
func GetFirefoxDirs(profileDir string) (error, string, string) {
	err1, dataDir := GetDataDir(profileDir)
	err2, cacheDir := GetCacheDir(profileDir)
	if err1 == nil && err2 == nil {
		return nil, dataDir, cacheDir
	} else {
		jerr := errors.Join(err1, err2)
		return fmt.Errorf("%s %q %w", ErrFirefoxCleaner, profileDir, jerr), dataDir, cacheDir
	}
}

func GetRootDataDir() (error, string) {
	appDataRoaming, err1 := os.UserConfigDir()
	if err1 != nil {
		return err1, ""
	}
	appDataLocal, err2 := os.UserCacheDir()
	if err2 != nil {
		return err2, ""
	}

	subPath := filepath.Join("Mozilla", "Firefox")
	if err, location := locate([]string{appDataRoaming, appDataLocal}, subPath, "ROOT"); err != nil {
		return err, ""
	} else {
		return nil, location
	}
}

// Profile-specific Cache directory
// *Windows: C:\Documents and Settings\<user>\Local Settings\Application Data\Mozilla\Firefox\Profiles\<profile>
// %APPDATA%\Mozilla\Firefox\Profiles\PROFILE\{cache2|Cache}
func GetCacheDir(profileDir string) (error, string) { // TODO: (Windows) needs to be verified!
	appDataRoaming, err1 := os.UserConfigDir()
	if err1 != nil {
		return err1, ""
	}
	appDataLocal, err2 := os.UserCacheDir()
	if err2 != nil {
		return err2, ""
	}

	subPath := filepath.Join("Mozilla", "Firefox", "Profiles", profileDir, "cache2")
	if err, location := locate([]string{appDataRoaming, appDataLocal}, subPath, "CACHE"); err != nil {
		return err, ""
	} else {
		return nil, location
	}
}

// Profile-specific Profile directory
// *Windows: C:\Documents and Settings\<user>\Application Data\Mozilla\Firefox\Profiles\<profile>
// %APPDATA%\Mozilla\Firefox\Profiles\7wdvp18g.default
func GetDataDir(profileDir string) (error, string) { // TODO: (Windows) needs to be verified!
	appDataRoaming, err1 := os.UserConfigDir()
	if err1 != nil {
		return err1, ""
	}
	appDataLocal, err2 := os.UserCacheDir()
	if err2 != nil {
		return err2, ""
	}

	subPath := filepath.Join("Mozilla", "Firefox", "Profiles", profileDir)
	if err, location := locate([]string{appDataRoaming, appDataLocal}, subPath, "DATA"); err != nil {
		return err, ""
	} else {
		return nil, location
	}
}

// attempt to find subPath in any of alternatives.
func locate(alternatives []string, subPath, prefix string) (error, string) {
	for _, lpath := range alternatives {
		name := filepath.Join(lpath, subPath)
		if cmn.IsDirectory(name) {
			return nil, name
		}
	}

	return fmt.Errorf("%s-ERR Couldn't find %q in any of %v", prefix, subPath, alternatives), ""
}
