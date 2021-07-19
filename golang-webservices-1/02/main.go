package main

import (
	"fmt"
	"io"
	"os"
)

const (
	Decorator = "├───"
	Tab = "\t"
)

type ObjectType struct {
	name string
	size int64
	isDir bool
	children []ObjectType
}

func worker(path string) (objects []ObjectType,  err error) {
	var (
		files []os.DirEntry
	)

	objects = make([]ObjectType, 0)

	files, err = os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		var (
			children []ObjectType
			fileInfo os.FileInfo
		)

		if file.IsDir() {
			children, err = worker(fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), file.Name()))
			if err != nil {
				return nil, err
			}

			objects = append(objects, ObjectType{name: file.Name(), isDir: true, children: children})
		} else {
			fileInfo, err = os.Stat(fmt.Sprintf("%s/%s", path, file.Name()))
			if err != nil {
				return nil, err
			}

			objects = append(objects, ObjectType{name: file.Name(), size: fileInfo.Size()})
		}
	}

	return objects, nil
}

func output(out io.Writer, printFiles bool, objects []ObjectType, level int) (err error) {
	for _, object := range objects {
		if !printFiles && !object.isDir {
			continue
		}

		for i := 0; i < level; i++ {
			_, err = out.Write([]byte(Tab))
			if err != nil {
				return err
			}
		}

		_, err = out.Write([]byte(fmt.Sprintf("%s%s", Decorator, object.name)))
		if err != nil {
			return err
		}

		if !object.isDir {
			_, err = out.Write([]byte(fmt.Sprintf(" (%db)", object.size)))
			if err != nil {
				return err
			}
		}

		_, err = out.Write([]byte(fmt.Sprintf("\n")))
		if err != nil {
			return err
		}

		if object.isDir {
			err = output(out, printFiles, object.children, level + 1)
		}
	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	var (
		objects []ObjectType
	)

	objects, err = worker(path)
	err = output(out, printFiles, objects, 0)

	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}