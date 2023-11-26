package rcs

import (
    "testing"
    "os"
)

const (
    txt1 = "a\nb\nc\n"
    txt2 = "a\nc\n"
    txt3 = "a\nb\nc\nd\n"
)

func TestDiff(t *testing.T) {
    d, err := diff(txt1, txt2)
    if err != nil {
        t.Fatalf("diff failed: %v", err)
    }
    t2, err := patch(txt1, d)
    if err != nil {
        t.Fatalf("patch failed: %v", err)
    }

    if txt2 != t2 {
        t.Errorf("Invalid result: %s != %s", t2, txt2)
    }
}

func TestBasic(t *testing.T) {
    defer os.Remove("rcs_test.txt")
    rcs, err := NewRCSFile("rcs_test.txt")
    if err != nil {
        t.Fatalf("can't create RCSFile: %v", err)
    }

    h1, err := rcs.Put(txt1, "comment1")
    if err != nil {
        t.Fatalf("commit failed: %v", err)
    }

    h2, err := rcs.Put(txt2, "comment2")
    if err != nil {
        t.Fatalf("commit failed: %v", err)
    }

    h3, err := rcs.Put(txt3, "comment3")
    if err != nil {
        t.Fatalf("commit failed: %v", err)
    }

    t1, err := rcs.GetVersion(h1)
    if err != nil {
        t.Fatalf("checkout failed: %v", err)
    }
    if t1 != txt1 {
        t.Errorf("Different texts: %s != %s", t1, txt1)
    }

    t2, err := rcs.GetVersion(h2)
    if err != nil {
        t.Fatalf("checkout failed: %v", err)
    }
    if t2 != txt2 {
        t.Errorf("Different texts: %s != %s", t2, txt2)
    }

    t3, err := rcs.GetVersion(h3)
    if err != nil {
        t.Fatalf("checkout failed: %v", err)
    }
    if t3 != txt3 {
        t.Errorf("Different texts: %s != %s", t3, txt3)
    }

    txt, err := rcs.Get()
    if err != nil {
        t.Fatalf("checkout failed: %v", err)
    }
    if txt != txt3 {
        t.Errorf("Different texts: %s != %s", txt, txt3)
    }
}
