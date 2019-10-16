package main

import (
	"errors"
	"fmt"
	stdioutil "io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	"gopkg.in/src-d/go-git.v4/utils/ioutil"
)

func main() {
	repo, err := openFromPath("/Users/alexis/go/gta")
	if err != nil {
		log.Fatal("can't open from path: ", err)
	}

	iter, err := repo.Branches()
	if err != nil {
		log.Fatal("can't get branches: ", err)
	}

	_ = iter.ForEach(func(reference *plumbing.Reference) error {
		fmt.Println(reference.Name())
		return nil
	})
}


func openFromPath(path string) (*git.Repository, error) {
	dot, wt, err := dotGitToOSFilesystems(path, true)
	if err != nil {
		return nil, err
	}

	if _, err := dot.Stat(""); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("repository .git does not exist")
		}

		return nil, err
	}

	s := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())

	return git.Open(s, wt)
}

func dotGitToOSFilesystems(path string, detect bool) (dot, wt billy.Filesystem, err error) {
	if path, err = filepath.Abs(path); err != nil {
		return nil, nil, err
	}
	var fs billy.Filesystem
	var fi os.FileInfo
	for {
		fs = osfs.New(path)
		fi, err = fs.Stat(".git")
		if err == nil {
			// no error; stop
			break
		}
		if !os.IsNotExist(err) {
			// unknown error; stop
			return nil, nil, err
		}
		if detect {
			// try its parent as long as we haven't reached
			// the root dir
			if dir := filepath.Dir(path); dir != path {
				path = dir
				continue
			}
		}
		// not detecting via parent dirs and the dir does not exist;
		// stop
		return fs, nil, nil
	}

	if fi.IsDir() {
		dot, err = fs.Chroot(".git")
		return dot, fs, err
	}

	dot, err = dotGitFileToOSFilesystem(path, fs)
	if err != nil {
		return nil, nil, err
	}

	return dot, fs, nil
}

func dotGitFileToOSFilesystem(path string, fs billy.Filesystem) (bfs billy.Filesystem, err error) {
	f, err := fs.Open(".git")
	if err != nil {
		return nil, err
	}
	defer ioutil.CheckClose(f, &err)

	b, err := stdioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	line := string(b)
	const prefix = "gitdir: "
	if !strings.HasPrefix(line, prefix) {
		return nil, fmt.Errorf(".git file has no %s prefix", prefix)
	}

	gitdir := strings.Split(line[len(prefix):], "\n")[0]
	gitdir = strings.TrimSpace(gitdir)
	if filepath.IsAbs(gitdir) {
		return osfs.New(gitdir), nil
	}

	return osfs.New(fs.Join(path, gitdir)), nil
}