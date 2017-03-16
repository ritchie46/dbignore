package main

import (
	"path/filepath"
	"fmt"
	"os"
	"os/exec"
	"github.com/fsnotify/fsnotify"

)



func main() {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		println(err)
	}

	defer watcher.Close()

	done := make(chan bool)
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

	err = watcher.Add("/home/ritchie46")

	<- done
	// execute("dropbox", "exclude")
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

