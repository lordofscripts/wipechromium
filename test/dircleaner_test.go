/* -----------------------------------------------------------------
 *				C o r a l y s   T e c h n o l o g i e s
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *						U n i t   T e s t
 *-----------------------------------------------------------------*/
package test

import (
	"fmt"
	"testing"
	"path/filepath"

	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"

	"lordofscripts/wipechromium"
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
		FSO{true, "Cache"},
		FSO{true, "Cache/Cache_Data"},
		FSO{false, "Cache/Cache_Data/index"},

		FSO{true, "Code Cache"},
		FSO{true, "Code Cache/js"},
		FSO{false, "Code Cache/js/index"},
		FSO{true, "Code Cache/wasm"},
		FSO{false, "Code Cache/wasm/index"},
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
		FSO{false, "Affiliation Database"},
		FSO{true, "AutofillStrikeDatabase"},
		FSO{false, "Bookmarks"},
		FSO{false, "Bookmarks.bak"},
		FSO{false, "BrowsingTopicsSiteData"},
		FSO{true, "BudgetDatabase"},
		FSO{true, "BudgetDatabase/LOCK"},
		FSO{false, "BudgetDatabase/LOG"},
		FSO{false, "BudgetDatabase/LOG.old"},
		FSO{true, "chrome_cart_db"},
		FSO{true, "ClientCertificates"},
		FSO{true, "commerce_subscription_db"},
		FSO{false, "Cookies"},
		FSO{false, "Cookies-journal"},

		FSO{true, "Extension Rules"},
		FSO{false, "Extension Rules/000003.log"},
		FSO{false, "Extension Rules/CURRENT"},
		FSO{false, "Extension Rules/LOCK"},
		FSO{false, "Extension Rules/LOG"},
		FSO{false, "Extension Rules/MANIFEST-000002"},

		FSO{true, "Extension Scripts"},
		FSO{false, "Extension Scripts/000003.log"},
		FSO{false, "Extension Scripts/CURRENT"},
		FSO{false, "Extension Scripts/LOCK"},
		FSO{false, "Extension Scripts/LOG"},
		FSO{false, "Extension Scripts/MANIFEST-000002"},

		FSO{true, "Extension State"},
		FSO{false, "Extension State/000003.log"},
		FSO{false, "Extension State/CURRENT"},
		FSO{false, "Extension State/LOCK"},
		FSO{false, "Extension State/LOG"},
		FSO{false, "Extension State/MANIFEST-000002"},

		FSO{false, "Favicons"},
		FSO{true, "File System"},
		FSO{true, "File System/001"},
		FSO{true, "File System/002"},
		FSO{true, "File System/Origins"},
		FSO{false, "History"},

		FSO{true, "Local Storage"},
		FSO{true, "Local Storage/leveldb"},

		FSO{false, "LOCK"},
		FSO{false, "LOG"},
		FSO{false, "Login Data"},
		FSO{false, "Login Data For Account"},
		FSO{false, "LOG.old"},
		FSO{false, "Preferences"},
		FSO{false, "PreferredApps"},
		FSO{false, "Secure Preferences"},
		FSO{true, "Sessions"},
		FSO{true, "Session Storage"},
		FSO{true, "Sync Data"},

		FSO{true, "Web Applications"},
		FSO{true, "Web Applications/Manifest Resources"},
		FSO{true, "Web Applications/Manifest Resources/nolmkcfonidpkniogdbnhmnnaepcehlc"},
		FSO{true, "Web Applications/Temp"},

		FSO{true, "Web Storage"},
		FSO{true, "Web Storage/10"},
		FSO{true, "Web Storage/20"},
		FSO{false, "Web Storage/QuotaManager"},
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
	path string
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
		showVFS(mfs, CacheDir)
	}

	cleaner := wipechromium.NewDirCleanerVFS(mfs, CacheDir, logx)
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
		showVFS(mfs, DataDir)
	}

	cleaner := wipechromium.NewDirCleanerVFS(mfs, DataDir, logx)
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

/* ----------------------------------------------------------------
 *					H e l p e r   F u n c t i o n s
 *-----------------------------------------------------------------*/

func createDummyTree(fs vfs.Filesystem, objects []FSO, root string) (int64, error) {
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
		} else {
			_, err = vfs.Create(fs, apath)
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
		return int64(len(entries)),err
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

func showVFS(fs vfs.Filesystem, root string) {
	fmt.Printf("Dumping FileSystem %q\n", root)

	if entries, err := fs.ReadDir(root); err == nil {
		for _, fso := range entries {
			label := "File"
			if fso.IsDir() {
				label = "Dir"
			}
			fmt.Printf("\t%5s %s\n", label, fso.Name())
		}
	} else {
		fmt.Println("VFS CANNOT SHOW", err.Error())
	}
}