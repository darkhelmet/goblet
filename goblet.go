package goblet

import (
    "archive/tar"
    "bitbucket.org/kardianos/osext"
    "bytes"
    "crypto/sha1"
    "fmt"
    "github.com/darkhelmet/goblet/elf"
    "io"
    "time"
)

type Asset struct {
    Name         string
    Data         []byte
    LastModified time.Time
    Sha1         string
}

func (a *Asset) Size() int {
    return len(a.Data)
}

func (a *Asset) Reader() *bytes.Reader {
    return bytes.NewReader(a.Data)
}

type Goblet struct {
    Files map[string]*Asset
}

func (g *Goblet) Get(name string) *Asset {
    return g.Files[name]
}

func Load() (*Goblet, error) {
    binary, err := osext.Executable()
    if err != nil {
        return nil, err
    }

    data, err := elf.ExtractSection(binary)
    if err != nil {
        return nil, err
    }

    r := tar.NewReader(bytes.NewReader(data))
    files := make(map[string]*Asset)
    for {
        hdr, err := r.Next()
        if err != nil {
            if err == io.EOF {
                break
            }
            return nil, fmt.Errorf("failed to get next entry: %s", err)
        }

        buffer := make([]byte, hdr.Size)
        _, err = io.ReadFull(r, buffer)
        if err != nil {
            return nil, fmt.Errorf("failed reading from archive %s: %s", hdr.Name, err)
        }

        files[hdr.Name] = &Asset{
            Name:         hdr.Name,
            Data:         buffer,
            LastModified: hdr.ModTime,
            Sha1:         fmt.Sprintf("%x", sha1.Sum(buffer)),
        }
    }

    return &Goblet{files}, nil
}
