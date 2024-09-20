/* -----------------------------------------------------------------
 *				C o r a l y s   T e c h n o l o g i e s
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *						U n i t   T e s t
 *-----------------------------------------------------------------*/
package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/lordofscripts/vfs"
	"github.com/lordofscripts/vfs/memfs"

	"github.com/lordofscripts/wipechromium"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	LOGGING_ENABLED bool = false
)

var (
	logx = wipechromium.NewConditionalLogger(LOGGING_ENABLED, "Test")

	ExceptionsCache []string = []string{}

	ChromiumCache = []FSO{
		{true, "Cache"},
		{true, "Cache/Cache_Data"},
		{false, "Cache/Cache_Data/index"},

		{true, "Code Cache"},
		{true, "Code Cache/js"},
		{false, "Code Cache/js/index"},
		{true, "Code Cache/wasm"},
		{false, "Code Cache/wasm/index"},
	}

	ExceptionsData []string = []string{
		"Bookmarks",
		"Extension Rules",
		"Extension State",
		"Extension Scripts",
		"File System",
		"Local Storage",
		"LOCK",
		"Preferences",
		"Secure Preferences",
		"Web Applications",
	}

	ChromiumData = []FSO{
		{false, "Affiliation Database"},
		{true, "AutofillStrikeDatabase"},
		{false, "Bookmarks"},
		{false, "Bookmarks.bak"},
		{false, "BrowsingTopicsSiteData"},
		{true, "BudgetDatabase"},
		{true, "BudgetDatabase/LOCK"},
		{false, "BudgetDatabase/LOG"},
		{false, "BudgetDatabase/LOG.old"},
		{true, "chrome_cart_db"},
		{true, "ClientCertificates"},
		{true, "commerce_subscription_db"},
		{false, "Cookies"},
		{false, "Cookies-journal"},

		{true, "Extension Rules"},
		{false, "Extension Rules/000003.log"},
		{false, "Extension Rules/CURRENT"},
		{false, "Extension Rules/LOCK"},
		{false, "Extension Rules/LOG"},
		{false, "Extension Rules/MANIFEST-000002"},

		{true, "Extension Scripts"},
		{false, "Extension Scripts/000003.log"},
		{false, "Extension Scripts/CURRENT"},
		{false, "Extension Scripts/LOCK"},
		{false, "Extension Scripts/LOG"},
		{false, "Extension Scripts/MANIFEST-000002"},

		{true, "Extension State"},
		{false, "Extension State/000003.log"},
		{false, "Extension State/CURRENT"},
		{false, "Extension State/LOCK"},
		{false, "Extension State/LOG"},
		{false, "Extension State/MANIFEST-000002"},

		{false, "Favicons"},
		{true, "File System"},
		{true, "File System/001"},
		{true, "File System/002"},
		{true, "File System/Origins"},
		{false, "History"},

		{true, "Local Storage"},
		{true, "Local Storage/leveldb"},

		{false, "LOCK"},
		{false, "LOG"},
		{false, "Login Data"},
		{false, "Login Data For Account"},
		{false, "LOG.old"},
		{false, "Preferences"},
		{false, "PreferredApps"},
		{false, "Secure Preferences"},
		{true, "Sessions"},
		{true, "Session Storage"},
		{true, "Sync Data"},

		{true, "Web Applications"},
		{true, "Web Applications/Manifest Resources"},
		{true, "Web Applications/Manifest Resources/nolmkcfonidpkniogdbnhmnnaepcehlc"},
		{true, "Web Applications/Temp"},

		{true, "Web Storage"},
		{true, "Web Storage/10"},
		{true, "Web Storage/20"},
		{false, "Web Storage/QuotaManager"},
	}
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type FSO struct {
	isDir bool
	path  string
}

/* ----------------------------------------------------------------
 *				U n i t  T e s t   F u n c t i o n s
 *-----------------------------------------------------------------*/
func Test_DirCleanerCache(t *testing.T) {
	const (
		CacheDir = "/home/pi/.cache/chromium/Profile 1"
	)
	mfs := memfs.Create()
	totalObjects := int64(0)

	if qty, err := createDummyTree(mfs, ChromiumCache, CacheDir); err != nil {
		t.Errorf(err.Error())
	} else {
		totalObjects = qty
		showVFS(mfs, CacheDir, false)
	}

	cleaner := wipechromium.NewDirCleanerVFS(mfs, CacheDir, wipechromium.SizeModeSI, logx)
	if err := cleaner.CleanUp(ExceptionsCache); err == nil {
		fmt.Printf("Cache cleaned: %d\n%s\n", cleaner.CleanedSize(), cleaner)
		if count := tallyTree(mfs, CacheDir); count != -1 && count < totalObjects {
			fmt.Println("Appears OK")
		} else {
			t.Errorf("Incomplete Cache deletion were %d now %d", totalObjects, count)
		}
	} else {
		t.Errorf(err.Error())
	}
}

func Test_DirCleanerData(t *testing.T) {
	const (
		DataDir = "/home/pi/.config/chromium/Profile 1"
	)
	mfs := memfs.Create()
	totalObjects := int64(0)

	if qty, err := createDummyTree(mfs, ChromiumData, DataDir); err != nil {
		t.Errorf(err.Error())
	} else {
		totalObjects = qty
		showVFS(mfs, DataDir, false)
	}

	cleaner := wipechromium.NewDirCleanerVFS(mfs, DataDir, wipechromium.SizeModeSI, logx)
	if err := cleaner.CleanUp(ExceptionsData); err == nil {
		fmt.Printf("Data cleaned: %d\n%s\n", cleaner.CleanedSize(), cleaner)
		if count := tallyTree(mfs, DataDir); count != -1 && count < totalObjects {
			fmt.Println("Appears OK")
		} else {
			t.Errorf("Incomplete Data deletion were %d now %d", totalObjects, count)
		}
	} else {
		t.Errorf(err.Error())
	}
}

func Test_MimicFileSystem(t *testing.T) {
	wd, _ := os.Getwd()
	upPath := filepath.Join(wd, "../")

	fmt.Printf("Mimicking %q as Memory VFS\n", upPath)
	mfs := memfs.Create()
	dirCnt, fileCnt, err := wipechromium.MimicFileSystem(upPath, mfs)
	if err != nil {
		t.Errorf("Couldn't mimic %v", err)
	} else {
		dC, fC := showVFS(mfs, upPath, false)
		if dirCnt != dC {
			t.Errorf("Expected %d dirs in VFS got %d", dirCnt, dC)
		}
		if fileCnt != fC {
			t.Errorf("Expected %d files in VFS got %d", fileCnt, fC)
		}

		fmt.Printf("OK. Mimicked %d dirs & %d files in MemVFS\n", dirCnt, fileCnt)
	}
}

/* ----------------------------------------------------------------
 *					H e l p e r   F u n c t i o n s
 *-----------------------------------------------------------------*/

func createDummyTree(fs vfs.Filesystem, objects []FSO, root string) (int64, error) {
	writeDummyFile := func(path, content string) {
		if err := vfs.WriteFile(fs, path, []byte(content), 0644); err != nil {
			fmt.Println("VFS.F.ERR", err)
		}
	}
	// (a) create the MemFS root
	if err := vfs.MkdirAll(fs, root, 0700); err != nil {
		return 0, err
	}

	// (b) create dummy objects
	for _, fso := range objects {
		var err error
		apath := filepath.Join(root, fso.path)
		if fso.isDir {
			err = fs.Mkdir(apath, 0700)
			dfilename := filepath.Join(apath, "dummy.txt")
			writeDummyFile(dfilename, "Anything")
			if IsFile(fs, dfilename) != wipechromium.Yes {
				fmt.Println("VFS.C.ERR where is ", dfilename)
			}
		} else {
			writeDummyFile(apath, "Test File")
		}

		if err != nil {
			fmt.Println("VFS ", err)
		}
	}

	if entries, err := fs.ReadDir(root); err == nil {
		fmt.Printf("Created %d root objects  on %q\n", len(entries), root)
		return int64(len(entries)), nil
	} else {
		fmt.Println("VFS ", err)
		return int64(len(entries)), err
	}
}

func tallyTree(fs vfs.Filesystem, root string) int64 {
	if entries, err := fs.ReadDir(root); err == nil {
		return int64(len(entries))
	} else {
		fmt.Println("VFS ", err.Error())
	}

	return -1
}

// recursively show a Virtual File System starting at root. The main
// caller should set isRecursing to false.
// Returns: number of directories & files encountered
func showVFS(fs vfs.Filesystem, root string, isRecursing bool) (dirCnt, fileCnt int64) {
	if !isRecursing {
		fmt.Printf("Dumping FileSystem %q\n", root)
	} else {
		fmt.Printf("\tDumping Subdir %q\n", root)
	}

	fileCnt = int64(0)
	dirCnt = int64(0)
	otherDirs := make([]string, 0)
	if entries, err := fs.ReadDir(root); err == nil {
		for _, fso := range entries {
			label := "File"
			if fso.IsDir() {
				label = "Dir"
				otherDirs = append(otherDirs, filepath.Join(root, fso.Name()))
				dirCnt++
			} else {
				fileCnt++
			}
			fmt.Printf("\t%5s %6d %s\n", label, fso.Size(), fso.Name())
		}

		// recurse directories
		for _, dirName := range otherDirs {
			subDirCnt, subFileCnt := showVFS(fs, dirName, true)
			dirCnt += subDirCnt
			fileCnt += subFileCnt
		}
	} else {
		fmt.Println("VFS CANNOT SHOW", err.Error())
	}

	return dirCnt, fileCnt
}

func IsFile(fs vfs.Filesystem, path string) wipechromium.TriState {
	if finfo, err := fs.Stat(path); err == nil {
		if !finfo.IsDir() {
			return wipechromium.Yes
		} else {
			return wipechromium.No
		}
	}
	// not exist or not a directory
	return wipechromium.Undecided
}
