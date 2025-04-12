# Go File Operations: A Comprehensive Guide

This guide explores file operations in Go, covering basic and advanced techniques for reading, writing, and manipulating files. Using packages like `os`, `io`, `path/filepath`, and `bufio`, Go provides robust tools for file handling. Learn to build efficient, secure, and idiomatic applications with this comprehensive resource.

## Table of Contents

- [Go File Operations: A Comprehensive Guide](#go-file-operations-a-comprehensive-guide)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
    - [Why File Operations Matter](#why-file-operations-matter)
  - [Basic File Operations](#basic-file-operations)
    - [Reading Files](#reading-files)
    - [Writing Files](#writing-files)
  - [Working with Paths](#working-with-paths)
    - [Path Manipulation](#path-manipulation)
    - [File Metadata](#file-metadata)
  - [Directory Operations](#directory-operations)
    - [Creating Directories](#creating-directories)
    - [Listing and Walking](#listing-and-walking)
  - [File Manipulation](#file-manipulation)
    - [Moving and Deleting](#moving-and-deleting)
    - [Permissions](#permissions)
  - [Advanced File Operations](#advanced-file-operations)
    - [File Locking](#file-locking)
    - [Memory-Mapped Files](#memory-mapped-files)
  - [Testing File Operations](#testing-file-operations)
    - [Mocking Files](#mocking-files)
    - [Temporary Files](#temporary-files)
  - [Best Practices](#best-practices)
    - [Error Handling](#error-handling)
    - [Resource Management](#resource-management)
    - [Security](#security)
  - [Common Patterns](#common-patterns)
    - [File Processing Pipeline](#file-processing-pipeline)
    - [Concurrent File Copy](#concurrent-file-copy)
  - [Key Takeaways](#key-takeaways)

---

## Introduction

File operations are essential for many Go applications, enabling data persistence, configuration management, and more. Go's standard library provides powerful, cross-platform tools for file handling, designed for safety and efficiency.

### Why File Operations Matter

- **Persistence**: Store data across sessions.
- **Configuration**: Manage settings and preferences.
- **Logging**: Record events and errors.
- **Processing**: Handle datasets for analysis or transformation.
- **Interoperability**: Exchange data with other systems.

---

## Basic File Operations

Go's `os`, `io`, and `bufio` packages form the backbone of file operations.

### Reading Files

**Example:**

```go
package main

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "os"
)

func main() {
    // Open file
    file, err := os.Open("data.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // Read entire file
    data, err := io.ReadAll(file)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Full content:", string(data))

    // Reset file pointer
    file.Seek(0, 0)

    // Read line by line
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        fmt.Println("Line:", scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}
```

**Key Points:**

- Use `os.Open` for read-only access.
- `io.ReadAll` reads entire files; `bufio.Scanner` is line-oriented.
- Always check errors and defer `Close`.

### Writing Files

**Example:**

```go
package main

import (
    "bufio"
    "log"
    "os"
)

func main() {
    // Create or overwrite file
    file, err := os.Create("output.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // Direct write
    data := []byte("Hello, World!\n")
    if _, err := file.Write(data); err != nil {
        log.Fatal(err)
    }

    // Buffered write
    writer := bufio.NewWriter(file)
    if _, err := writer.WriteString("Buffered write\n"); err != nil {
        log.Fatal(err)
    }
    if err := writer.Flush(); err != nil {
        log.Fatal(err)
    }
}
```

**Key Points:**

- `os.Create` truncates or creates files.
- Use `bufio.Writer` for efficient writes.
- Call `Flush` to ensure buffered data is written.

---

## Working with Paths

The `path/filepath` package ensures cross-platform path handling.

### Path Manipulation

**Example:**

```go
package main

import (
    "fmt"
    "path/filepath"
)

func main() {
    // Join paths
    path := filepath.Join("dir", "subdir", "file.txt")
    fmt.Println("Path:", path) // dir/subdir/file.txt

    // Extract components
    fmt.Println("Dir:", filepath.Dir(path))      // dir/subdir
    fmt.Println("Base:", filepath.Base(path))    // file.txt
    fmt.Println("Ext:", filepath.Ext(path))      // .txt

    // Clean path
    messy := "dir/../dir/./file.txt"
    fmt.Println("Clean:", filepath.Clean(messy)) // dir/file.txt
}
```

**Key Functions:**

- `Join`: Combines path components.
- `Dir`, `Base`, `Ext`: Extract parts.
- `Clean`: Normalizes paths, resolving `..` and `.`.

### File Metadata

**Example:**

```go
package main

import (
    "fmt"
    "log"
    "os"
)

func main() {
    info, err := os.Stat("data.txt")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Name:", info.Name())        // data.txt
    fmt.Println("Size:", info.Size())        // File size in bytes
    fmt.Println("IsDir:", info.IsDir())      // false
    fmt.Println("ModTime:", info.ModTime())  // Last modified time
    fmt.Println("Mode:", info.Mode())        // File permissions
}
```

**Use Cases:**

- Validate file existence with `os.Stat`.
- Check metadata like size or permissions.
- Handle `os.IsNotExist(err)` for missing files.

---

## Directory Operations

Go simplifies directory creation and traversal.

### Creating Directories

**Example:**

```go
package main

import (
    "log"
    "os"
)

func main() {
    // Create single directory
    if err := os.Mkdir("mydir", 0755); err != nil {
        log.Fatal(err)
    }

    // Create directory tree
    if err := os.MkdirAll("parent/child/grandchild", 0755); err != nil {
        log.Fatal(err)
    }
}
```

**Key Points:**

- `Mkdir` creates one directory; `MkdirAll` creates a tree.
- Use standard permissions (e.g., `0755` for directories).

### Listing and Walking

**Example:**

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
)

func main() {
    // List directory contents
    entries, err := os.ReadDir("mydir")
    if err != nil {
        log.Fatal(err)
    }
    for _, entry := range entries {
        fmt.Println(entry.Name(), entry.IsDir())
    }

    // Walk directory tree
    err = filepath.WalkDir("parent", func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }
        fmt.Printf("Visited: %s (Dir: %v)\n", path, d.IsDir())
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

**Key Points:**

- `os.ReadDir` lists a single directory.
- `filepath.WalkDir` traverses recursively, with fine-grained control.

---

## File Manipulation

### Moving and Deleting

**Example:**

```go
package main

import (
    "log"
    "os"
)

func main() {
    // Rename/move file
    if err := os.Rename("old.txt", "new.txt"); err != nil {
        log.Fatal(err)
    }

    // Delete file
    if err := os.Remove("new.txt"); err != nil {
        log.Fatal(err)
    }

    // Delete directory and contents
    if err := os.RemoveAll("mydir"); err != nil {
        log.Fatal(err)
    }
}
```

**Key Points:**

- `Rename` moves or renames files atomically.
- `Remove` deletes files; `RemoveAll` deletes directories recursively.

### Permissions

**Example:**

```go
package main

import (
    "log"
    "os"
)

func main() {
    // Set file permissions
    if err := os.Chmod("data.txt", 0644); err != nil {
        log.Fatal(err)
    }

    // Change ownership (Unix only)
    if err := os.Chown("data.txt", os.Getuid(), os.Getgid()); err != nil {
        log.Fatal(err)
    }

    // Open with specific permissions
    file, err := os.OpenFile("data.txt", os.O_RDWR, 0600)
    if err != nil {
        log.Fatal(err)
    }
    file.Close()
}
```

**Key Points:**

- Use `Chmod` for Unix-style permissions (e.g., `0644` for files).
- `Chown` is platform-specific (Unix).
- Set permissions explicitly when creating files.

---

## Advanced File Operations

### File Locking

File locking ensures exclusive access, useful for concurrent processes.

**Example:**

```go
package main

import (
    "log"
    "os"
    "syscall"
)

type FileLocker struct {
    file *os.File
}

func NewFileLocker(filename string) (*FileLocker, error) {
    f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        return nil, err
    }
    return &FileLocker{file: f}, nil
}

func (l *FileLocker) Lock() error {
    return syscall.Flock(int(l.file.Fd()), syscall.LOCK_EX)
}

func (l *FileLocker) Unlock() error {
    return syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
}

func (l *FileLocker) Close() error {
    return l.file.Close()
}

func main() {
    locker, err := NewFileLocker("lock.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer locker.Close()

    if err := locker.Lock(); err != nil {
        log.Fatal(err)
    }
    defer locker.Unlock()

    // Perform exclusive operations
    if _, err := locker.file.WriteString("Locked write\n"); err != nil {
        log.Fatal(err)
    }
}
```

**Key Points:**

- Use `syscall.Flock` for advisory locking (Unix-specific).
- Ensure proper cleanup with `Close` and `Unlock`.
- Suitable for coordinating multiple processes.

### Memory-Mapped Files

Memory mapping boosts performance for large files.

**Example:**

```go
package main

import (
    "log"
    "os"
    "syscall"
)

func mapFile(filename string) ([]byte, error) {
    f, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    stat, err := f.Stat()
    if err != nil {
        return nil, err
    }

    data, err := syscall.Mmap(int(f.Fd()), 0, int(stat.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
    if err != nil {
        return nil, err
    }
    return data, nil
}

func main() {
    data, err := mapFile("data.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer syscall.Munmap(data)

    log.Println("Mapped content:", string(data))
}
```

**Key Points:**

- `syscall.Mmap` maps files to memory (Unix-specific).
- Ideal for large datasets or frequent access.
- Requires `Munmap` to release resources.

---

## Testing File Operations

### Mocking Files

Use `io.Reader`/`io.Writer` interfaces to mock file operations.

**Example:**

```go
package main

import (
    "bytes"
    "io"
    "testing"
)

func processFile(r io.Reader) (string, error) {
    data, err := io.ReadAll(r)
    if err != nil {
        return "", err
    }
    return string(data), nil
}

func TestProcessFile(t *testing.T) {
    input := bytes.NewReader([]byte("test data"))
    got, err := processFile(input)
    if err != nil {
        t.Fatal(err)
    }
    want := "test data"
    if got != want {
        t.Errorf("Got %q, want %q", got, want)
    }
}
```

### Temporary Files

Use `os.CreateTemp` for safe testing.

**Example:**

```go
package main

import (
    "os"
    "testing"
)

func TestWriteFile(t *testing.T) {
    tmp, err := os.CreateTemp("", "test*.txt")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmp.Name())

    if _, err := tmp.WriteString("test"); err != nil {
        t.Fatal(err)
    }

    data, err := os.ReadFile(tmp.Name())
    if err != nil {
        t.Fatal(err)
    }
    if string(data) != "test" {
        t.Errorf("Got %q, want %q", data, "test")
    }
}
```

**Key Points:**

- Mock with interfaces for portability.
- Use temporary files to avoid side effects.
- Clean up test files with `defer`.

---

## Best Practices

### Error Handling

- **Always Check Errors**: Use `if err != nil` for every operation.
- **Wrap Errors**: Add context with `fmt.Errorf("...: %w", err)`.
- **Handle Edge Cases**: Account for missing files, permissions, and EOF.

### Resource Management

- **Defer Closes**: Ensure files are closed with `defer file.Close()`.
- **Flush Buffers**: Call `Flush` for `bufio.Writer`.
- **Stop Operations**: Use `os.Remove` or `os.RemoveAll` for cleanup.

### Security

- **Set Permissions**: Use `0644` for files, `0755` for directories.
- **Validate Paths**: Prevent directory traversal with `filepath.Clean`.
- **Lock Files**: Use locking for concurrent access when needed.

**Example (Safe Write):**

```go
func safeWrite(filename, content string) error {
    file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        return fmt.Errorf("open %s: %w", filename, err)
    }
    defer file.Close()

    if _, err := file.WriteString(content); err != nil {
        return fmt.Errorf("write %s: %w", filename, err)
    }
    return nil
}
```

---

## Common Patterns

### File Processing Pipeline

Chain transformations on file data.

**Example:**

```go
package main

import (
    "fmt"
    "io"
    "os"
    "strings"
)

type FileProcessor struct {
    input  string
    output string
    stages []func([]byte) []byte
}

func (p *FileProcessor) AddStage(stage func([]byte) []byte) {
    p.stages = append(p.stages, stage)
}

func (p *FileProcessor) Process() error {
    data, err := os.ReadFile(p.input)
    if err != nil {
        return fmt.Errorf("read %s: %w", p.input, err)
    }

    for _, stage := range p.stages {
        data = stage(data)
    }

    if err := os.WriteFile(p.output, data, 0644); err != nil {
        return fmt.Errorf("write %s: %w", p.output, err)
    }
    return nil
}

func main() {
    p := FileProcessor{input: "in.txt", output: "out.txt"}
    p.AddStage(func(b []byte) []byte { return []byte(strings.ToUpper(string(b))) })
    p.AddStage(func(b []byte) []byte { return append(b, "! "...) })
    if err := p.Process(); err != nil {
        fmt.Println("Error:", err)
    }
}
```

### Concurrent File Copy

Copy multiple files concurrently with workers.

**Example:**

```go
package main

import (
    "fmt"
    "io"
    "os"
    "sync"
)

func copyFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    if _, err := io.Copy(out, in); err != nil {
        return err
    }
    return out.Close()
}

func copyFiles(files map[string]string) error {
    var wg sync.WaitGroup
    errors := make(chan error, len(files))

    for src, dst := range files {
        wg.Add(1)
        go func(src, dst string) {
            defer wg.Done()
            if err := copyFile(src, dst); err != nil {
                errors <- fmt.Errorf("copy %s to %s: %w", src, dst, err)
            }
        }(src, dst)
    }

    wg.Wait()
    close(errors)

    for err := range errors {
        if err != nil {
            return err
        }
    }
    return nil
}

func main() {
    files := map[string]string{
        "file1.txt": "copy1.txt",
        "file2.txt": "copy2.txt",
    }
    if err := copyFiles(files); err != nil {
        fmt.Println("Error:", err)
    }
}
```

---

## Key Takeaways

- **Basics**: Use `os.Open`, `os.Create`, and `bufio` for efficient I/O.
- **Paths**: Leverage `path/filepath` for cross-platform compatibility.
- **Directories**: Manage with `MkdirAll`, `ReadDir`, and `WalkDir`.
- **Manipulation**: Handle moves, deletes, and permissions carefully.
- **Advanced**: Use locking and memory mapping for special cases.
- **Testing**: Mock interfaces and use temporary files.
- **Best Practices**:
  - Always defer resource cleanup.
  - Check and wrap errors with context.
  - Secure files with proper permissions and path validation.

This guide equips you to handle file operations in Go with confidence, from simple reads to complex concurrent pipelines.

_Last Updated: April 11, 2025_
