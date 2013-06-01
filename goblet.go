package goblet

import (
    "archive/tar"
    "bytes"
    "fmt"
    "github.com/darkhelmet/goblet/elf"
    "io"
    "os"
)

type Goblet struct {
    Files map[string][]byte
}

func (g *Goblet) Get(name string) []byte {
    return g.Files[name]
}

func (g *Goblet) GetReader(name string) io.Reader {
    return bytes.NewReader(g.Get(name))
}

func Load() (*Goblet, error) {
    data, err := elf.ExtractSection(os.Args[0])
    if err != nil {
        return nil, err
    }

    var buffer bytes.Buffer
    r := tar.NewReader(bytes.NewReader(data))
    files := make(map[string][]byte)
    for {
        hdr, err := r.Next()
        if err != nil {
            if err == io.EOF {
                break
            }
            return nil, fmt.Errorf("failed to get next entry: %s", err)
        }
        buffer.Reset()
        _, err = io.Copy(&buffer, r)
        if err != nil {
            return nil, fmt.Errorf("failed reading from archive %s: %s", hdr.Name, err)
        }

        files[hdr.Name] = buffer.Bytes()
    }

    return &Goblet{files}, nil
}
