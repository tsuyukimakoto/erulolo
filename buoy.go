package main

/*
  cat /proc/sys/fs/inotify/max_user_watches
	8192
	ubuntu14.04 can watch 8192 files.
*/

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"golang.org/x/exp/inotify"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

type Config struct {
	Buoy BuoySetting
}

type BuoySetting struct {
	IgnoreExtensions []string "toml:ignoreExtensions"
	IgnoreDotDirectory bool "toml:ignoreDotDirectory"
	Targets []string "toml:targets"
}

func isIgnore(_path string, ignore_exts []string) bool {
	for _, _ext := range ignore_exts {
		if _ext == strings.ToLower(path.Ext(_path)) {
			return true
		}
	}
	return false
}

func absDirectory(root string, path string) string {
	_rel, _ := filepath.Rel(root, path)
	_directory, _ := filepath.Abs(_rel)
	return _directory
}

func recursiveDirectories(root string, ignoreDotDirectory bool) []string {
	targets := make([]string, 0, 100)

	evfunc := func(path string, info os.FileInfo, err error) error {
		if info != nil && info.IsDir() {
			// TODO more clear
			if ignoreDotDirectory == true && strings.Contains(path, "/.") == true {
				return nil
			} else {
				// targets = append(targets, absDirectory(root, path))
				targets = append(targets, path)
			}
		}
		return err
	}

	err := filepath.Walk(root, evfunc)
	if err != nil {
		fmt.Println(1, err)
	}
	return targets
}

func addWatch(watcher *inotify.Watcher, path string) {
	watcher.AddWatch(path, syscall.IN_CREATE)
	watcher.AddWatch(path, syscall.IN_DELETE)
	watcher.AddWatch(path, syscall.IN_CLOSE_WRITE)
	log.Println("Watch: ", path)
}

func Watch() {
	var config Config
	_, err := toml.DecodeFile("config.tml", &config)
	if err != nil {
		panic(err)
	}
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("config.Buoy.Targets:", config.Buoy.Targets)
	log.Println("config.Buoy.IgnoreExtensions:", config.Buoy.IgnoreExtensions)
	log.Println("config.Buoy.IgnoreDotDirectory:", config.Buoy.IgnoreDotDirectory)
	for _, targets := range config.Buoy.Targets {
		addWatch(watcher, targets)
		for _, target := range recursiveDirectories(targets, config.Buoy.IgnoreDotDirectory) {
			// TODO make watching event type clear.
			addWatch(watcher, target)
		}
	}
	for {
		select {
		case ev := <-watcher.Event:
			switch ev.Mask {
			// TODO IN_CREATE comes and file type is directory, AddWatch it.
			// TODO IN_DELETE comes and file type is directory, RemoveWatch
			case syscall.IN_CLOSE_WRITE, syscall.IN_CREATE, syscall.IN_DELETE:
				if isIgnore(ev.Name, config.Buoy.IgnoreExtensions) != true {
					// TODO notify this in some way.
					log.Println("event:", ev)
				}
			}
		case err := <-watcher.Error:
			log.Println("error:", err)
		}
	}
}
