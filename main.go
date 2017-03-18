package main

import (
	"path/filepath"
	"fmt"
	"os"
	"os/exec"
	"github.com/fsnotify/fsnotify"
	"runtime"
	//"io/ioutil"

	"bufio"
	//"reflect"

)

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
	fmt.Println(ignore)

	watcher := newWatcher()
	_ = watcher

	dirs := Walker(ignore, root)
	fmt.Println(dirs)

	// Add a watcher to all directories.
	var dirMap = make(map[string]bool, len(dirs))
	for _, v := range dirs {
		dirMap[v] = true
		watcher.Add(v)
	}

	done := make(chan bool)
	<- done

	//d, _ := define_directories(root)
	//
	//a := len(d)
	//fmt.Println(a)
	//

	//err := watcher.Add("C:/Users/vik/Dropbox/Code")
	//
	//if err != nil {
	//	fmt.Println(err)
	//}

	//done := make(chan bool)
	//<- done
	// execute("dropbox", "exclude")
}

// Copied form path module. Only this slice will not be sorted.
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
func Walker(ignore []string, root string) []string {
	var dirs = []string{root}

	// Retrieve all names in the root directory
	for _, d := range readDirNames(root) {
		if d.IsDir() && d.Name() != ".git"{
			var walk_dir = true

			// If the directory is in the ignore slice, the directory may not be walked.
			for _, ignore_dir := range ignore {
				if ignore_dir == d.Name() {
					walk_dir = false
					break
				}
			}
			if walk_dir {
				dirs = append(dirs, Walker(ignore, filepath.Join(root, d.Name()))...)
			}
		}
	}
	return dirs
}


// Run a shell/ cmd command.
func execute(bin string, args ...string) {
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

// Create a new file watcher. This function describes the callback events.
func newWatcher() *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		println(err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				eventHandler(event, watcher)

				println("Event:", event.Name, event.Op.String())
			case err := <-watcher.Errors:
				println("Errors:", err)
			}
		}
	}()
	return watcher
}

// Handle the file events
func eventHandler(e fsnotify.Event, w *fsnotify.Watcher) {
	action := e.Op.String()
	if action == "REMOVE" || action == "RENAME"{
		// fsnotify cannot remove watcher on removed files/ dirs.
		// w.Remove(e.Name) <- Throws a panic.
	} else {
		fs, err := os.Stat(e.Name)

		if err != nil {
			fmt.Println(err)
		}

		if fs.IsDir() {
			if e.Op.String() == "CREATE" {
				fmt.Println("Added ", e.Name)
				w.Add(e.Name)
			}

		}
	}
}