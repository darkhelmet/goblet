package elf

import (
    "debug/elf"
    "errors"
    "github.com/darkhelmet/goblet"
)

func HasSection(path string) bool {
    file, err := elf.Open(path)
    if err != nil {
        return false
    }
    defer file.Close()

    return file.Section(goblet.SectionName) != nil
}

func ExtractSection(path string) ([]byte, error) {
    file, err := elf.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    section := file.Section(goblet.SectionName)
    if section == nil {
        return nil, errors.New("no goblet section found")
    }

    return section.Data()
}
