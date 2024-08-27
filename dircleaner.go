/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2024 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package wipechromium

import (
	"fmt"
	"os"
	"path/filepath"
	"slices" // GO v1.18
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	SizeModeStd SizeMode = iota // full numeric size
	SizeModeSI                  // Sistema Internacional: 1KB = 1000
	SizeModeIEC                 // Binary System: 1KB = 1024
)

var _ IDirCleaner = (*DirCleaner)(nil)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

type IDirCleaner interface {
	String() string
	CleanUp(exceptions []string) error
	CleanedSize() int64
}

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type DirCleaner struct {
	Root        string
	cleanedSize int64
	removedQty  int
	skippedQty  int
	logx        ILogger
}

type SizeMode uint

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewDirCleaner(root string, logger ...ILogger) *DirCleaner {
	const cName = "DirCleaner"
	var logCtx ILogger
	if len(logger) == 0 {
		logCtx = NewConditionalLogger(false, cName)
	} else {
		logCtx = logger[0].InheritAs(cName)
	}
	return &DirCleaner{root, 0, 0, 0, logCtx}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (s SizeMode) String() string {
	switch s {
	case SizeModeStd:
		return "Standard"
	case SizeModeSI:
		return "International"
	case SizeModeIEC:
		return "Binary"
	default:
		panic("Unknown size mode")
	}
}

func (d *DirCleaner) String() string {
	return fmt.Sprintf("DirCleaner %q del:%d skip:%d size:%s", d.Root,
		d.removedQty,
		d.skippedQty,
		AddThousands(d.cleanedSize, ','))
}

func (d *DirCleaner) CleanUp(exceptions []string) error {
	d.cleanedSize = 0
	d.removedQty = 0
	d.skippedQty = 0
	entries, err := os.ReadDir(d.Root)
	if err != nil {
		d.logx.Print(err)
		return err
	}

	var executor func(string) error = nil

	execRemoveRecursive := func(path string) error {
		return os.RemoveAll(path)
	}

	execRemoveSingle := func(path string) error {
		return os.Remove(path)
	}

	for _, item := range entries {
		if !slices.Contains(exceptions, item.Name()) {
			// count up
			size := int64(0)
			fullPath := filepath.Join(d.Root, item.Name())
			// get file/dir size
			if finfo, err := item.Info(); err == nil {
				if finfo.IsDir() {
					executor = execRemoveRecursive
					size, _ = GetDirectorySize(fullPath)
					d.logx.Printf("%8d D %s", size, fullPath)
				} else {
					executor = execRemoveSingle
					size = finfo.Size()
					d.logx.Printf("%8d F %s", size, fullPath)
				}
			} else {
				d.logx.Print("DirCleaner WARN:", err)
			}

			// delete
			if err := executor(fullPath); err != nil {
				return err
			} else {
				d.cleanedSize += size
			}
			d.removedQty += 1
		} else {
			d.skippedQty += 1
			d.logx.Printf("DirCleaner skipping %s", item.Name())
		}
	}
	return nil
}

func (d *DirCleaner) CleanedSize() int64 {
	return d.cleanedSize
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

func GetDirectorySize(folder string) (int64, error) {
	// Step 1: remember subdirectories we must recurse
	var folders []string
	// Step 2: read directory and handle errors.
	dirRead, err := os.Open(folder)
	if err != nil {
		return 0, err
	}
	dirFiles, err := dirRead.Readdir(0)
	if err != nil {
		return 0, err
	}

	// Step 3: sum up Size of all files in the directory.
	sum := int64(0)
	for _, fileHere := range dirFiles {
		if fileHere.IsDir() {
			// Size() returns the size of the directory entry not the actual
			// sum of file sizes in that directory. Therefore, we skip it for later.
			folders = append(folders, filepath.Join(folder, fileHere.Name()))
		} else {
			csize := fileHere.Size()
			sum += csize
			//fmt.Printf("%8d %5t %s\n", csize, fileHere.IsDir(), fileHere.Name())
		}
	}

	// Step 4: close directory and return the sum.
	dirRead.Close()

	// Step: 5: Iterate recursively through the directories we encountered
	for _, dn := range folders {
		folderSize, _ := GetDirectorySize(dn)
		sum += folderSize
	}
	return sum, nil
}
