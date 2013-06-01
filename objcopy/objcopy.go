package objcopy

import (
    "flag"
    "fmt"
    "github.com/darkhelmet/goblet"
    "os/exec"
)

var objcopy = flag.String("objcopy", "objcopy", "The path to objcopy")

func Gobletize(input, archive string) error {
    cmd := exec.Command(*objcopy, "--add-section", fmt.Sprintf("%s=%s", goblet.SectionName, archive), input)
    return cmd.Run()
}
