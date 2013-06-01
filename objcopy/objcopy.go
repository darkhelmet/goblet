package objcopy

import (
    "flag"
    "fmt"
    "github.com/darkhelmet/goblet"
    "os/exec"
)

var objcopy = flag.String("objcopy", "objcopy", "The path to objcopy")

func RemoveGoblet(input string) error {
    cmd := exec.Command(*objcopy, "--remove-section", goblet.SectionName, input)
    return cmd.Run()
}

func Gobletize(input, archive string) error {
    cmd := exec.Command(*objcopy, "--add-section", fmt.Sprintf("%s=%s", goblet.SectionName, archive), input)
    return cmd.Run()
}
