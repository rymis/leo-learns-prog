// Simple revision control system implementation.

package rcs

import (
	"crypto"
    "fmt"
	"encoding/json"
	"encoding/hex"
	"os"
	"sync"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var globalLock sync.Mutex
const emptyHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

// RCSFile allows to read/write file with versioning
type RCSFile struct {
    // Path to file
    Path string
}

type VersionInfo struct {
    Version string    `json:"version"`
    Time    float64   `json:"time"`
    Comment string    `json:"comment"`
    Parent  string    `json:"parent"`
}

type rcsHist struct {
    Version string `json:"version"`
    Time float64   `time:"date"`
    Comment string `json:"comment"`
    Data    string `json:"data"`
    Parent  string `json:"parent"`
}

type rcsFile struct {
    Current rcsHist `json:"current"`
    History map[string]rcsHist `json:"history"`
}

// Create new RCSFile
func NewRCSFile(filename string) (*RCSFile, error) {
    // TODO: check if file exists
    return &RCSFile{filename}, nil   
}

func (rcs *RCSFile) load() (*rcsFile, error) {
    f, err := os.Open(rcs.Path)
    if err != nil { // Create empty file
        res := &rcsFile{}
        res.Current.Data = ""
        res.Current.Time = ftime()
        res.Current.Version = emptyHash
        res.History = make(map[string]rcsHist)
        return res, nil
    }
    defer f.Close()

    loader := json.NewDecoder(f)
    res := &rcsFile{}
    err = loader.Decode(res)
    if err != nil {
        return nil, err
    }

    err = res.CheckIntegrity()
    if err != nil {
        return nil, err
    }

    return res, nil
}

func (rcs *RCSFile) save(r *rcsFile) error {
    tmpName := rcs.Path + "~"
    f, err := os.Create(tmpName)
    if err != nil {
        return err
    }

    encoder := json.NewEncoder(f)
    encoder.SetIndent("", "  ")
    err = encoder.Encode(r)
    if err != nil {
        f.Close()
        os.Remove(tmpName)
        return err
    }

    err = f.Close()
    if err != nil {
        os.Remove(tmpName)
        return err
    }

    // This operation is atomic so we minimize risk of loosing some data
    return os.Rename(tmpName, rcs.Path)
}

func (rcs *RCSFile) Get() (string, error) {
    globalLock.Lock()
    defer globalLock.Unlock()

    r, err := rcs.load()
    if err != nil {
        return "", err
    }

    return r.Current.Data, nil
}

func (rcs *RCSFile) Put(data string, comment string) (string, error) {
    // This version goes to history and we create the new one
    hist := rcsHist{}
    hist.Data = data
    hist.Comment = comment
    hist.Time = ftime()

    globalLock.Lock()
    defer globalLock.Unlock()

    r, err := rcs.load()
    if err != nil {
        return "", err
    }
    hist.Version = hashVersion(data, len(r.History))

    hist.Parent = r.Current.Version
    r.Current.Data, err = diff(data, r.Current.Data)
    if err != nil {
        return "", err
    }
    r.History[r.Current.Version] = r.Current
    r.Current = hist

    return hist.Version, rcs.save(r)
}

func (rcs* RCSFile) Versions() ([]VersionInfo, error) {
    globalLock.Lock()
    r, err := rcs.load()
    globalLock.Unlock()

    if err != nil {
        return nil, err
    }

    res := make([]VersionInfo, 0, len(r.History) + 1)
    res = append(res, r.Current.ToInfo())
    ver := r.Current.Parent
    for ver != "" {
        hitem := r.History[ver]
        res = append(res, hitem.ToInfo())
        ver = hitem.Parent
    }

    return res, nil
}

func (rcs* RCSFile) GetVersion(ver string) (string, error) {
    globalLock.Lock()
    r, err := rcs.load()
    globalLock.Unlock()
    if err != nil {
        return "", err
    }

    h := r.Current
    txt := h.Data
    for h.Parent != "" {
        if h.Version == ver {
            return txt, nil
        }

        h = r.History[h.Parent]
        txt, err = patch(txt, h.Data)
        if err != nil {
            return "", err
        }
    }

    return "", fmt.Errorf("Unknown version: %s", ver)
}

func (rcs *rcsFile) CheckIntegrity() error {
    // TODO: check file integrity
    return nil
}

func ftime() float64 {
    return float64(time.Now().UnixMilli()) / 1000.0
}

func timef(f float64) time.Time {
    t := int64(f * 1000.0)
    return time.UnixMilli(t)
}

func diff(txt1, txt2 string) (string, error) {
    dmp := diffmatchpatch.New()
    d := dmp.DiffMain(txt1, txt2, true)
    return dmp.DiffToDelta(d), nil
}

func patch(txt, diff string) (string, error) {
    dmp := diffmatchpatch.New()
    d, err := dmp.DiffFromDelta(txt, diff)
    if err != nil {
        return "", err
    }
    return dmp.DiffText2(d), nil
}

func hashVersion(content string, count int) string {
    h := crypto.SHA256.New()
    fmt.Fprintf(h, "%s", content)
    val := h.Sum(nil)
    return fmt.Sprintf("%s:%d", hex.EncodeToString(val), count)
}

func (h rcsHist) ToInfo() (res VersionInfo) {
    res.Comment = h.Comment
    res.Version = h.Version
    res.Time = h.Time
    res.Parent = h.Parent
    return
}

