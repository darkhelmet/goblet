package main

import (
    "archive/tar"
    "bytes"
    "flag"
    "fmt"
    "github.com/darkhelmet/goblet/elf"
    "github.com/darkhelmet/goblet/objcopy"
    "io"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
)

var (
    input     = flag.String("input", "", "The input file to use")
    extract   = flag.String("extract", "", "Extract the archive to the given file")
    gobletize = flag.String("gobletize", "", "The directory to pack")
)

func walker(base string, tw *tar.Writer) filepath.WalkFunc {
    return func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() {
            hdr, err := tar.FileInfoHeader(info, "")
            if err != nil {
                return fmt.Errorf("failed building header: %s", err)
            }

            rel, err := filepath.Rel(base, path)
            if err != nil {
                return err
            }
            hdr.Name = rel

            err = tw.WriteHeader(hdr)
            if err != nil {
                return fmt.Errorf("failed writing header: %s", err)
            }

            file, err := os.Open(path)
            if err != nil {
                return fmt.Errorf("failed opening %s: %s", path, err)
            }

            _, err = io.Copy(tw, file)
            if err != nil {
                return fmt.Errorf("failed adding %s to archive: %s", path, err)
            }
        }

        return nil
    }
}

func archive(path string) string {
    arc, err := ioutil.TempFile(os.TempDir(), "gobletize")
    if err != nil {
        log.Fatalf("failed making temp file: %s", err)
    }
    defer arc.Close()

    t := tar.NewWriter(arc)
    defer t.Close()

    err = filepath.Walk(path, walker(path, t))
    if err != nil {
        log.Fatalf("failed walking %s: %s", path, err)
    }

    return arc.Name()
}

func Gobletize() {
    err := objcopy.Gobletize(*input, archive(*gobletize))
    if err != nil {
        log.Fatalf("gobletize failed: %s", err)
    }
}

func ExtractArchive() {
    data, err := elf.ExtractSection(*input)
    if err != nil {
        log.Fatalf("failed extracting section from %s: %s", *input, err)
    }

    file, err := os.OpenFile(*extract, os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        log.Fatalf("failed to open %s for writing: %s", *extract, err)
    }
    defer file.Close()

    _, err = io.Copy(file, bytes.NewReader(data))
    if err != nil {
        log.Fatalf("failed copy: %s", err)
    }
}

func main() {
    flag.Parse()

    if *input == "" {
        log.Fatalf("no input file specified, use -input <path>")
    }

    if *gobletize != "" {
        Gobletize()
    } else if *extract != "" {
        ExtractArchive()
    } else {
        log.Fatalf("provided either -gobletize or -extract to gobletize or extract")
    }
}
