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
	fmt.Println(ignore)

	watcher := newWatcher()
	_ = watcher

	dirs := Walker(ignore, root)

	for _, v := range dirs {
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

// Walk the root directories and follow every dir that is not ignored.
func Walker(ignore []string, root string) []string {
	//fmt.Println("new Walker call", root)

	var dirs = []string{root}
	for _, d := range readDirNames(root) {
		if d.IsDir() {
			var walk_dir = true
			var ignore_dir string
			for _, ignore_dir = range ignore {
				if ignore_dir == d.Name() {
					walk_dir = false
					break
				}
			}
			if walk_dir {
				//w.Add(p)
				dirs = append(dirs, Walker(ignore, filepath.Join(root, d.Name()))...)
			}
		}
	}
	return dirs
}


// Find al the directories recursively
// Path is the root
func define_directories(path string) ([]string, error){
	directories := [] string{}

	// Callback function from filepath.Walk
	var wf = func(path string, fi os.FileInfo, err error) error {

		d := filepath.Dir(path)

		if len(directories) == 0 {
			directories = append(directories, d)

		} else if d != directories[len(directories) - 1] {
			directories = append(directories, d)
		}
		return err
	}

	err := filepath.Walk(path, wf)

	if err != nil {
		return nil, err
	}
	return directories, nil
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

	//defer watcher.Close()


	go func() {
		for {
			select {
			case event := <-watcher.Events:
				println("Event:", event.Name)
			case err := <-watcher.Errors:
				println("Errors:", err)
			}
		}
	}()
	return watcher
}