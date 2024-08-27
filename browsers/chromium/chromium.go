/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Privacy & Junk cleaner for Chromium-based browsers
 *-----------------------------------------------------------------*/
package chromium

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cmn "lordofscripts/wipechromium"
	"lordofscripts/wipechromium/browsers"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	PERMS                 = 0700
	RecreateCacheDir bool = false
)

// ensure we qualify as supported browser plugin
var (
	_ browsers.IBrowsers = (*ChromiumCleaner)(nil)

	// Don't delete these on Chromium root, i.e. ~/.config/chromium/
	ChromiumExceptions []string = []string{ // top-level exceptions (add user profiles)
		"Avatars",
		"Default",              // default profile
		"extensions_crx_cache", // just in case
		"System Profile",
	}
	// Don't delete these on Profile root, i.e. ~/.config/chromium/{profile name}/
	ProfileExceptions []string = []string{ // user profile exceptions
		"Bookmarks",
		"Bookmarks.bak",
		"LOCK",
		"Preferences",
		"PreferredApps",
		// {profile_name}/directories
		"Extension Rules",
		"Extensions",
		"Extension Scripts",
		"Extension State",
		"File System", // Progressive Web Apps keep their data here! i.e. Novelist!
		"Local Extension Settings",
		"Web Applications", // therein remove Temp
	}
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type ChromiumCleaner struct {
	Class       browsers.Browser
	ProfileName string
	CacheRoot   string
	ProfileRoot string
	cleanedSize int64
	sizeMode    cmn.SizeMode
	logx        cmn.ILogger
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewChromiumCleaner(profile string, smode cmn.SizeMode, logger ...cmn.ILogger) *ChromiumCleaner {
	const cName = "ChromiumCleaner"
	var logCtx cmn.ILogger
	if len(logger) == 0 {
		logCtx = cmn.NewConditionalLogger(false, cName)
	} else {
		logCtx = logger[0].InheritAs(cName)
	}

	ChromiumDataDir, ChromiumCachesDir := GetChromiumDirs()

	return &ChromiumCleaner{browsers.ChromiumBrowser,
		strings.Trim(profile, " \t"),
		filepath.Join(ChromiumCachesDir, profile),
		filepath.Join(ChromiumDataDir, profile),
		0,
		smode,
		logCtx,
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (c *ChromiumCleaner) String() string {
	var reportedSize string = ""
	switch c.sizeMode {
	case cmn.SizeModeSI:
		reportedSize = cmn.ByteCountSI(c.cleanedSize)
		break
	case cmn.SizeModeIEC:
		reportedSize = cmn.ByteCountIEC(c.cleanedSize)
		break
	case cmn.SizeModeStd:
		fallthrough
	default:
		reportedSize = cmn.AddThousands(c.cleanedSize, ',')
		break
	}

	c.logx.Print("Cleaned ", cmn.AddThousands(c.cleanedSize, ','))
	c.logx.Print("Cleaned ", cmn.ByteCountSI(c.cleanedSize))
	c.logx.Print("Cleaned ", cmn.ByteCountIEC(c.cleanedSize))
	return fmt.Sprintf("%sCleaner %q cleaned %s", c.Class, c.ProfileName, reportedSize)
}

func (c *ChromiumCleaner) Name() browsers.Browser {
	return browsers.ChromiumBrowser
}

// Top level function to clear a Chromium user profile directory. Rather than
// saving important data to a Temp directory and then restoring (as previous version)
// now we simply go through the top level with a list of exceptions
// Example: clearProfile("Profile 1")
func (c *ChromiumCleaner) ClearProfile(doCache, doProfile bool) (error, int) {
	fmt.Printf("Clearing profile %q\n", c.ProfileName)

	c.cleanedSize = 0
	if len(c.ProfileName) == 0 {
		// we can only operate in AppData root only as no profile is given
		return cmn.ErrNoProfile, 40
	}

	// 1. Profile Cache
	if doCache {
		if err := c.clearCache(); err != nil {
			return err, 50
		}
	}

	// 2. Profile Data
	if doProfile {
		if err := c.eraseProfile(); err != nil {
			return err, 60
		}

		if err := c.clearExtensions(); err != nil {
			return err, 70
		}
	}

	c.logx.Printf("Profile %q cleared of private/junk data", c.ProfileName)
	return nil, 0
}

// This function should be implemented in all wiper browser plugins.
// It should print out the supposed location of the Data & Cache directories
// so that the user can verify prior to running the program for the 1st time.
func (c *ChromiumCleaner) Tell() bool {
	ChromiumDataDir, ChromiumCachesDir := GetChromiumDirs()
	dataExists := cmn.IsDirectory(ChromiumDataDir)
	cacheExists := cmn.IsDirectory(ChromiumCachesDir)
	fmt.Println("Chromium Directories:")
	fmt.Println("\tData : ", ChromiumDataDir, dataExists)
	fmt.Println("\tCache: ", ChromiumCachesDir, cacheExists)
	return dataExists && cacheExists
}

// Browser data for ALL profiles. A user account has ONE browser AppDataRoot,
// but therein it may have more than one user Profile, each with its
// different settings, extensions, bookmarks, etc.
// Returns: true if GetDataDir() is the root of all user account profiles.
func (c *ChromiumCleaner) IdentifyAppDataRoot() bool {
	return IdentifyAppDataRoot()
}

// A user profile's cache directory.
// Returns: true if GetCacheDir()+profile is a valid browser Cache directory.
func (c *ChromiumCleaner) IdentifyProfileCache(profile string) bool {
	return IdentifyProfileCache(profile)
}

// A user profile specific data (extensions, Bookmarks, etc.).
// Applies to each of the profile directories under the directory identified
// by IdentifyAppDataRoot().
// Returns: true if directory contains browser user profile data & settings.
func (c *ChromiumCleaner) IdentifyProfileData(profile string) bool {
	return IdentifyProfileData(profile)
}

/* ----------------------------------------------------------------
 *					I n t e r n a l 	M e t h o d s
 *-----------------------------------------------------------------*/

// Clears the entire cache dir of a profile
// Example: clearCache("Profile 1")
func (c *ChromiumCleaner) clearCache() error {
	fmt.Println("\tClearing cache...")

	cacheSize, err := cmn.GetDirectorySize(c.CacheRoot)
	if err != nil {
		c.logx.Printf("clearCache WARN %s", err)
	}

	if !IdentifyProfileCache(c.ProfileName) {
		c.logx.Printf("%s: %s", cmn.ErrNotBrowserCache, c.CacheRoot)
		return cmn.ErrNotBrowserCache
	}

	// 'Cache' 'Code Cache' and sometimes 'Storage'
	if err := os.RemoveAll(c.CacheRoot); err != nil {
		return cmn.WrapError(err, 41, "Could not remove cache dir %q", c.CacheRoot)
	}

	if RecreateCacheDir {
		if err := os.Mkdir(c.CacheRoot, PERMS); err != nil {
			c.logx.Printf("Could not recreate CACHE %s", c.CacheRoot)
		}
	}

	c.cleanedSize += cacheSize
	fmt.Printf("\tDeleted %s bytes from cache\n", cmn.AddThousands(cacheSize, ','))
	c.logx.Print("clearCache DONE")
	return nil
}

// erases a User Profile but keeps important profile data such as
// extensions and settings.
func (c *ChromiumCleaner) eraseProfile() error {
	fmt.Println("\tClearing profile")

	// (a )Identify it is a profile directory
	if !IdentifyProfileData(c.ProfileName) {
		return cmn.ErrNotBrowserProfile
	}

	// (b) we are going to clean the profile's top level
	var filter cmn.IDirCleaner
	filter = cmn.NewDirCleaner(c.ProfileRoot)

	// (c) except these important profile items
	err := filter.CleanUp(ProfileExceptions)
	if err != nil {
		c.logx.Print(err)
		return cmn.WrapError(err, 51, "EraseProfile fault.")
	}

	c.cleanedSize += filter.CleanedSize()
	fmt.Printf("\t...Erased %d bytes\n", c.cleanedSize)
	return nil
}

func (c *ChromiumCleaner) clearExtensions() error {
	// (a) extension subdirs to cleanup
	categories := []string{
		"Extension Scripts",
		"Extension State",
		"Extension Rules",
	}

	// (b) iterate through profile extension category subdirs
	for _, subDir := range categories {
		fmt.Println("\tClearing ", subDir, "...")

		// (b.1) root of that extension data category
		root := filepath.Join(c.ProfileRoot, subDir)
		// (b.2) glob patterns to remove matching files
		patterns := []string{
			"*.log",
			"LOG*",
		}
		// (b.3) delete those files based on pattern matching
		if err := c.removeWithPatterns(root, patterns); err != nil {
			c.logx.Print("WARN", err)
		}
	}

	fmt.Println("\t...Cleared extension junk")
	return nil
}

// Removes all files matching a Pattern at Dir.
// Example: removeWithPattern("/home/lordofscripts/.cache", "*.log")
func (c *ChromiumCleaner) removeWithPatterns(dir string, patterns []string) error {
	for _, pattern := range patterns {
		glob := dir + string(os.PathSeparator) + pattern
		files, err := filepath.Glob(glob)
		if err != nil {
			return err
		}

		for _, fname := range files {
			var fileSize int64 = 0
			if finfo, err := os.Stat(fname); err != nil {
				return err
			} else {
				fileSize = finfo.Size()
			}
			// remove file or empty directory
			if err := os.Remove(fname); err != nil {
				return err
			}
			c.cleanedSize += fileSize
		}
	}

	return nil
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// Identify directory (from GetDataDir()) as a Chromium user account directory.
// Remember every user can have several profiles. This identifies just the
// root where the profiles are located
// At least 'System Profile', 'Avatars' & 'Safe Browsing' exist
func IdentifyAppDataRoot() bool {
	appdata := GetDataDir()
	return cmn.IsDirectory(filepath.Join(appdata, "System Profile")) &&
		cmn.IsDirectory(filepath.Join(appdata, "Default")) &&
		cmn.IsDirectory(filepath.Join(appdata, "Avatars")) &&
		cmn.IsDirectory(filepath.Join(appdata, "Safe Browsing"))
}

// Identify GetCacheDir() as a proper Chromium Cache directory
// Both 'Cache' & 'Code Cache' dirs exist
func IdentifyProfileCache(profile string) bool {
	cache := GetCacheDir()
	return cmn.IsDirectory(filepath.Join(cache, profile, "Cache")) &&
		cmn.IsDirectory(filepath.Join(cache, profile, "Code Cache"))
}

// At least Bookmarks exist && Cookies
func IdentifyProfileData(profile string) bool {
	user := filepath.Join(GetDataDir(), profile)
	return cmn.IsDirectory(filepath.Join(user, "Extension Rules")) &&
		cmn.IsFile(filepath.Join(user, "Preferences")) == cmn.Yes &&
		cmn.IsFile(filepath.Join(user, "Bookmarks")) == cmn.Yes
}
