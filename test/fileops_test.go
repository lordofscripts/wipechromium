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
	"reflect"
	"testing"

	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"

	cmn "github.com/lordofscripts/wipechromium"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				U n i t  T e s t   F u n c t i o n s
 *-----------------------------------------------------------------*/

// Tries constructor's default and the explicit setting (dry) and verifies
func Test_DryRunNOP(t *testing.T) {
	fmt.Println("➤ DryRun Mapping NOP")

	d := cmn.NewDryRunner()

	verifyFunc := func() {
		if err := d.AssertMapping(cmn.DryRunTargetNOP); err != nil {
			t.Errorf("Unexpected NOP mapping %v", err)
			outcome(false)
		} else {
			outcome(true)
		}
	}

	fmt.Println("\tChecking with Default")
	verifyFunc()

	fmt.Println("\tChecking with Explicit Enable()")
	d.Enable()
	verifyFunc()
}

// Configures it for real OS run (wet) and verifies.
func Test_DryRunOS(t *testing.T) {
	const MAPPING = cmn.DryRunTargetOS
	fmt.Println("➤ DryRun Mapping OS")

	d := cmn.NewDryRunner()
	d.Disable()

	if d.GetMode() != MAPPING {
		t.Errorf("Expected %s mode (wet mode)", MAPPING)
	}

	if err := d.AssertMapping(MAPPING); err != nil {
		t.Errorf("Unexpected %s mapping %v", MAPPING, err)
		outcome(false)
	} else {
		outcome(true)
	}
}

// Configures it for Virtual Filesystem (damp) and verifies
func Test_DryRunVFS(t *testing.T) {
	const MAPPING = cmn.DryRunTargetVFS
	fmt.Println("➤ DryRun Mapping VFS")

	mfs := memfs.Create()
	d := cmn.NewDryRunner()
	d.EnableOn(mfs)

	if d.GetMode() != MAPPING {
		t.Errorf("Expected %s mode", MAPPING)
	}

	if err := d.AssertMapping(MAPPING); err != nil {
		t.Errorf("Unexpected %s mapping %v", MAPPING, err)
		outcome(false)
	} else {
		outcome(true)
	}
}

// Creates a Memory-based Virtual Filesystem, then runs a series of filesystem
// operations first in dry mode (NOP) where nothing should change, and then
// in damp mode (VFS) where the virtual filesystem SHOULD get altered.
func Test_DryRunOnMemFS(t *testing.T) {
	fmt.Println("➤ DryRun Execute NOP & VFS")

	// A. Setup
	const (
		MY_ROOT   = "/tmp/home/toering"
		RECURSIVE = true
	)
	fmt.Println("  § Setup")

	// 1. Create a virtual memory file system
	mfs := memfs.Create()
	if mfs == nil {
		t.Errorf("Could not create MemFS")
		t.FailNow()
	}
	fmt.Printf("\tPath Separator: %s\n", string(mfs.PathSeparator()))

	// 2. Create a set of files & directories in that MemFS
	err := createFS(mfs, MY_ROOT)
	if err != nil {
		t.Errorf("Unable to create VFS: %v", err)
		t.FailNow()
	}

	// 2.1 count how many files & dirs were created therein (Initial state)
	dirCntO, fileCntO := countVFS(mfs, MY_ROOT, true)
	fmt.Printf("\t=  MemFS Original: %d Dirs, %d Files\n", dirCntO, fileCntO)

	// 3. Instantiate a Dry Runner to operate (later) on that MemFS
	dry := cmn.NewDryRunner() // default NOP

	// B. Test Cases
	type subCaseT struct {
		Mode   cmn.DryRunTarget
		Result bool
		Title  string
	}
	var subCases []subCaseT = []subCaseT{
		{cmn.DryRunTargetNOP, false, "I (NOP)"},
		{cmn.DryRunTargetVFS, false, "II (VFS)"},
	}

	// B.1 Run test sub cases
	var overallOutcome bool = true
	for _, subCase := range subCases {
		fmt.Printf("  § Test Case %s\n", subCase.Title)
		var dirDelta, fileDelta int64
		var expected int64
		subCase.Result = true

		// B.1.1 FSOps run mode configuration
		if subCase.Mode == cmn.DryRunTargetNOP {
			dry.Enable()
		} else if subCase.Mode == cmn.DryRunTargetVFS {
			dry.EnableOn(mfs)
			//dry.Dump()
		} else {
			panic("Test logic error") // in case someone adds something
		}

		// B.1.2 do filesystem operations
		dirDelta, fileDelta, err = doOperations(dry, MY_ROOT, subCase.Mode)
		if err != nil {
			t.Errorf("Operations failed: %v", err)
			subCase.Result = false
		}

		// B.1.2.1 For NOP we must ignore the Deltas
		if subCase.Mode == cmn.DryRunTargetNOP {
			dirDelta = 0
			fileDelta = 0
		}
		fmt.Printf("\tΔ Deltas: %+d Dirs, %+d Files\n", dirDelta, fileDelta) // \u0394

		// B.1.3 Verify the OS Action/Query Mappings are as expected
		if err := dry.AssertMapping(subCase.Mode); err != nil {
			t.Errorf("Unexpected %q mapping %v\n", subCase.Mode, err)
			subCase.Result = false
		}

		// B.1.4 outcome of DIR operations
		dirCntX, fileCntX := countVFS(mfs, MY_ROOT, RECURSIVE)
		fmt.Printf("\t= MemFS Now: %d Dirs, %d Files\n", dirCntX, fileCntX)

		// B.1.4.1 Check DIR ops
		expected = dirCntO + dirDelta
		if dirCntX != expected {
			t.Errorf("%s FileOps failed. Directory count %d != %d", subCase.Mode, dirCntX, expected)
			subCase.Result = false
		}

		// B.1.4.2 Check FILE operations
		expected = fileCntO + fileDelta
		if fileCntX != expected {
			t.Errorf("%s FileOps failed. File count %d != %d", subCase.Mode, fileCntX, expected)
			subCase.Result = false
		}

		// B.1.5 Subcase outcome
		outcome(subCase.Result)
		//dry.Dump()

		// B.2 Accumulate overall outcome
		overallOutcome = overallOutcome && subCase.Result
	}

	fmt.Println("  § Dump FS")
	dry.DumpFS(MY_ROOT)

	// C. Overall outcome of all subcases
	fmt.Println("  § Epilogue")
	outcome(overallOutcome)
}

/* ----------------------------------------------------------------
 *					H e l p e r   F u n c t i o n s
 *-----------------------------------------------------------------*/

// creates a dummy filesystem at 'root' with a set of files and directories
func createFS(mfs vfs.Filesystem, root string) error {
	type FSO struct {
		isDir bool
		path  string
	}

	writeDummyFile := func(path, content string) {
		if err := vfs.WriteFile(mfs, path, []byte(content), 0644); err != nil {
			fmt.Println("VFS.F.ERR", err)
		} else {
			fmt.Printf("\t\tℲ %s\n", path)
		}
	}

	FilesData := []FSO{
		{false, "root1.txt"},
		{false, "root2.pdf"},
		{false, "root3.doc"},
		{true, "Documents"},
		{true, "Documents/Confidential"},
		{false, "Documents/Confidential/secret.txt"},
		{true, "Pictures"},
		{false, "Pictures/image1.jpg"},
		{true, "Downloads"},
		{false, "Downloads/app1.tgz"},
		{false, "Downloads/app2.tgz"},
		{false, "Downloads/app3.tar"},
	}

	// (a) create the MemFS root
	if err := vfs.MkdirAll(mfs, root, 0700); err != nil {
		return err
	} else {
		fmt.Printf("\t\tⅮ %s\n", root)
	}

	// (b) create dummy objects
	for _, fso := range FilesData {
		var err error
		apath := filepath.Join(root, fso.path)
		if fso.isDir {
			err = mfs.Mkdir(apath, 0700)
			fmt.Printf("\t\tⅮ %s\n", apath)
			dfilename := filepath.Join(apath, "dummy.txt")
			writeDummyFile(dfilename, "Anything")
		} else {
			writeDummyFile(apath, "Test File")
		}

		if err != nil {
			fmt.Println("VFS ", err)
		}
	}

	return nil
}

// recursively show a Virtual File System starting at root. The main
// caller should set isRecursing to false.
// Returns: number of directories & files encountered
func countVFS(fs vfs.Filesystem, root string, isRecursing bool) (dirCnt, fileCnt int64) {
	fileCnt = int64(0)
	dirCnt = int64(0)
	otherDirs := make([]string, 0)
	if entries, err := fs.ReadDir(root); err == nil {
		for _, fso := range entries {
			if fso.IsDir() {
				otherDirs = append(otherDirs, filepath.Join(root, fso.Name()))
				dirCnt++
			} else {
				fileCnt++
			}
		}

		// recurse directories
		for _, dirName := range otherDirs {
			subDirCnt, subFileCnt := countVFS(fs, dirName, true)
			dirCnt += subDirCnt
			fileCnt += subFileCnt
		}
	} else {
		fmt.Println("VFS Couldn't count", err.Error())
	}

	return dirCnt, fileCnt
}

// helper to perform Action* operations on the selected filesystem
func doOperations(dry *cmn.DryRun, root string, mode cmn.DryRunTarget) (int64, int64, error) {
	fmt.Printf("\tOperating on %s\n", dry)

	var dirDelta, fileDelta int64
	// creates 3 directories under root
	err := dry.MkDirAll(filepath.Join(root, "A", "B", "C"), 0755) // NO +3, 0
	if err != nil {
		return 0, 0, fmt.Errorf("MkDirAll(%s): %v", mode, err)
	}
	dirDelta += 3

	// removes 1 file
	err = dry.Remove(filepath.Join(root, "Pictures/image1.jpg")) // NO 0, -1
	if err != nil {
		return 0, 0, fmt.Errorf("Remove(%s): %v", mode, err)
	}
	fileDelta -= 1

	// removes Documents & Documents/Confidential & 3 files
	err = dry.RemoveAll(filepath.Join(root, "Documents")) // NO -2, -2
	if err != nil {
		return 0, 0, fmt.Errorf("Remove(%s): %v", mode, err)
	}
	dirDelta -= 2
	fileDelta -= 3

	// rename
	oldFN := filepath.Join(root, "Downloads/app3.tar")
	newFN := filepath.Join(root, "Downloads/app3.zip")
	err = dry.Rename(oldFN, newFN)
	if err != nil {
		return 0, 0, fmt.Errorf("Rename(%s): %v", mode, err)
	}

	// extra check!
	if err := dry.AssertMapping(mode); err != nil {
		return 0, 0, fmt.Errorf("Unexpected %q mapping %v\n", mode, err)
	}

	return dirDelta, fileDelta, nil
}

// Works for Remove*
func compareA(func1, func2 func(string) error) bool {
	return reflect.ValueOf(func1).Pointer() == reflect.ValueOf(func2).Pointer()
}

// Works for Mkdir*
func compareB(func1, func2 func(string, os.FileMode) error) bool {
	return reflect.ValueOf(func1).Pointer() == reflect.ValueOf(func2).Pointer()
}

// Works for Rename
func compareC(func1, func2 func(string, string) error) bool {
	return reflect.ValueOf(func1).Pointer() == reflect.ValueOf(func2).Pointer()
}

// Works for QueryIsFile
func compareQ1(func1, func2 func(string) cmn.TriState) bool {
	return reflect.ValueOf(func1).Pointer() == reflect.ValueOf(func2).Pointer()
}

func outcome(ok bool) {
	if ok {
		fmt.Println("\t* ✔ OK")
	} else {
		fmt.Println("\t* ✘ FAILED")
	}
}
