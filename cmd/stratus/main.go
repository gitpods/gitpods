package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/oklog/oklog/pkg/group"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "stratus"

	app.Commands = []cli.Command{{
		Name:   "dev",
		Usage:  "Runs gitloud in development mode",
		Action: actionDev,
	}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func actionDev(c *cli.Context) error {
	var g group.Group
	g.Add(RunGitloud, func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	})

	g.Add(RunWebpack, func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	})

	return g.Run()
}

// RunGitloud runs a development server and restarts it with a new build if files change.
func RunGitloud() error {
	builds := make(chan bool)

	go BuildForever(builds)

	go func() {
		build()
		builds <- true
	}()

	var cmd *exec.Cmd
	for {
		<-builds
		if cmd != nil {
			cmd.Process.Kill()
		}
		cmd = exec.Command("./dist/gitloud", "web")
		go func() {
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Println(err)
				return
			}
		}()
	}

	return nil
}

// BuildForever watches the filesystem and builds a new binary if something changes.
// It notifies a channel that a build was created
func BuildForever(builds chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
		return
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op != fsnotify.Chmod && event.Name != "" {
					start := time.Now()
					if err := build(); err == nil { // only notify and log if binary was created successfully.
						log.Println("built a new binary in", time.Since(start))
						builds <- true // notify that a new build was successfully created
					}
					watcher.Remove(event.Name)
					watcher.Add(event.Name)
				}
			case err := <-watcher.Errors:
				log.Println(err)
			}
		}
	}()

	err = watcher.Add("./cmd/gitloud/main.go")
	if err != nil {
		log.Println(err)
		return
	}

	select {}
}

func build() error {
	cmd := exec.Command("go", "build", "-v", "-i", "-o", "./dist/gitloud", "./cmd/gitloud")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// RunWebpack starts webpack in the background and watches for changes
func RunWebpack() error {
	file := "./webpack.config.js"
	_, err := os.Stat(file)
	if err != nil {
		// webpack config not found
		return nil
	}

	cmd := exec.Command(filepath.Join("node_modules", ".bin", "webpack"), "--watch")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
