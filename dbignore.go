package main

import (
	"path/filepath"
	"fmt"
	"os"
	"os/exec"
	"github.com/fsnotify/fsnotify"
	"runtime"
	"bufio"
	"bytes"
	"strings"
)

type Ignore struct {
	ignore []string // Slice of the directory names that will be ignored.
	ignoreMap map[string]bool // Verify if dir in ignore. Maps the path to a bool.
	w *fsnotify.Watcher
}

// Return a slice with all the lines of a file
func readLine(path string) []string {
	f, _ := os.Open(path)
	defer f.Close()

	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)

	var s []string
	for scan.Scan() {
		s = append(s, scan.Text())
	}
	return s
}


func main() {
	var root string
	if runtime.GOOS == "windows" {
		root = "C:/Users/vik/Dropbox/Code"
	}else {
		root = "/home/ritchie46/Dropbox"
	}


	ignore := readLine(root + "/.dbignore")
	ignoreMap := make(map[string]bool, len(ignore))
	for _, v := range ignore {
		ignoreMap[v] = true
	}

	i := Ignore{ignore, ignoreMap, nil}


	fmt.Println(ignore)

	i.newWatcher()


	dirs := Walker(i, root)
	fmt.Println(dirs)

	// Add a watcher to all directories.
	var dirMap = make(map[string]bool, len(dirs))
	for _, v := range dirs {
		dirMap[v] = true
		i.w.Add(v)
	}



	done := make(chan bool)
	<- done
}

// Return the files and directories in a directory. The result is unsorted.
func readDirNames(dirname string) []os.FileInfo {
	f, err := os.Open(dirname)
	if err != nil {
		return nil
	}
	names, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil
	}
	return names
}

// Walk the root directories and follow every directory that is not ignored.
// .git directory is automatically ignored.
func Walker(i Ignore, root string) []string {
	var dirs = []string{root}

	// Retrieve all names in the root directory
	for _, d := range readDirNames(root) {
		if d.IsDir() && d.Name() != ".git"{
			var walk_dir = true

			// If the directory is in the ignore slice, the directory may not be walked.
			for _, ignore_dir := range i.ignore {
				if ignore_dir == d.Name() {
					walk_dir = false
					break
				}
			}
			if walk_dir {  // Walk this directory
				dirs = append(dirs, Walker(i, filepath.Join(root, d.Name()))...)
			} else { // ignore the directory in dropbox
				go dbexclude(filepath.Join(root, d.Name()))
			}
		}
	}
	return dirs
}


// Run a shell/ cmd command.
func execute(bin string, args ...string) string{
	cmd := exec.Command(bin, args...)
	var b  bytes.Buffer

	cmd.Stdout = &b
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return b.String()
}

// Create a new file watcher. This function describes the callback events.
func (i *Ignore)newWatcher() {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		println(err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				i.eventHandler(event)

				fmt.Println("Event:", event.Name, event.Op.String())
			case err := <-watcher.Errors:
				fmt.Println("Errors:", err)
			}
		}
	}()
	i.w = watcher
}

// Handle the file events
func (i *Ignore)eventHandler(e fsnotify.Event) {
	action := e.Op.String()
	if action == "REMOVE" || action == "RENAME"{
		// fsnotify cannot remove watcher on removed files/ dirs.
		// i.w.Remove(e.Name) <- Throws a panic.
	} else if e.Op.String() == "CREATE"{
		fs, err := os.Stat(e.Name)

		if err != nil {
			fmt.Println(err)
		}

		if fs.IsDir() {
			if i.ignoreMap[filepath.Base(e.Name)] {
				go dbexclude(e.Name)
			} else {
				fmt.Println("Added ", e.Name, "to watcher")
				i.w.Add(e.Name)
			}
		}
	}
}

func dbexclude(path string) {
	s := execute("dropbox", "exclude", "add", path)

	fmt.Println(s)


}

func dbinclude(im map[string]bool) {
	// list all ignored directories.
	ls := execute("dropbox", "exclude")
	temp := strings.Split(ls, "\n")
	for _, v := range temp {
		if v != "Excluded: " {
			abs_v, _ := filepath.Abs(v)

			if im[filepath.Base(abs_v)] {
				_, err := os.Stat(abs_v)
				if os.IsNotExist(err) {
					// include the files to sync
					s := execute("dropbox", "exclude", "remove", v)
					fmt.Println(s, abs_v)
				}
			}
		}
	}
}