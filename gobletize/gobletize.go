package main

import (
    "archive/tar"
    "flag"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "path/filepath"
)

var (
    objcopy = flag.String("objcopy", "objcopy", "The path to objcopy")
    input   = flag.String("input", "", "The input file to use")
    dump    = flag.Bool("dump", false, "Dump contents of existing section to stdout")
    dir     = flag.String("dir", "", "The directory to pack")
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

func main() {
    flag.Parse()

    if *input == "" {
        log.Fatalf("no input file specified, use -input <path>")
    }

    if *dir == "" {
        log.Fatalf("no dir specified, use -dir <dir>")
    }

    assets := archive(*dir)
    cmd := exec.Command(*objcopy, "--add-section", fmt.Sprintf("goblet=%s", assets), *input)
    if err := cmd.Run(); err != nil {
        log.Fatalf("objcopy failed: %s", err)
    }
}
