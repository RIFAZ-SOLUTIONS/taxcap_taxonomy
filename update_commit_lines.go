package main

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    root := "knowledge"
    err := filepath.Walk(root, processPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error walking the path %q: %v\n", root, err)
        os.Exit(1)
    }
}

func processPath(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    if info.IsDir() {
        return nil
    }
    if strings.ToLower(info.Name()) != "qna.yaml" {
        return nil
    }

    input, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("failed to open %s: %v", path, err)
    }
    defer input.Close()

    scanner := bufio.NewScanner(input)
    var lines []string

    for scanner.Scan() {
        line := scanner.Text()
        trimmed := strings.TrimLeft(line, " \t")
        if strings.HasPrefix(trimmed, "commit:") {
            // preserve indentation
            indent := line[:len(line)-len(trimmed)]
            line = indent + "commit: 3e10cd8"
        }
        lines = append(lines, line)
    }
    if err := scanner.Err(); err != nil {
        return fmt.Errorf("error reading %s: %v", path, err)
    }

    output := strings.Join(lines, "\n") + "\n"
    if err := ioutil.WriteFile(path, []byte(output), info.Mode()); err != nil {
        return fmt.Errorf("failed to write %s: %v", path, err)
    }

    fmt.Printf("Updated %s\n", path)
    return nil
}

