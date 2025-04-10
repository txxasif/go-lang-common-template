# Go File Operations: A Comprehensive Guide

## Table of Contents

1. [Introduction](#introduction)
2. [Basic File Operations](#basic-file-operations)
3. [Working with Paths](#working-with-paths)
4. [File Manipulation](#file-manipulation)
5. [Advanced File Operations](#advanced-file-operations)
6. [Best Practices](#best-practices)
7. [Common Patterns](#common-patterns)

## Introduction

File operations are fundamental to many Go applications. Understanding how to work with files efficiently and safely is crucial for building robust applications.

### Why File Operations Matter

1. **Data Persistence**: Store data between program runs
2. **Configuration**: Read and write application settings
3. **Logging**: Record application events and errors
4. **Data Processing**: Handle large datasets
5. **Interoperability**: Exchange data with other systems

## Basic File Operations

### Understanding File Operations in Go

Go provides a comprehensive set of tools for file operations through the `os` and `io` packages. These operations are designed to be safe, efficient, and easy to use.

#### Reading Files

```go
// Basic file reading
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

// Read line by line
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    // Process line
}
```

**Key Concepts:**

1. **File Handles**: Represent open files
2. **Buffering**: Improves performance
3. **Error Handling**: Essential for robustness
4. **Resource Cleanup**: Prevent resource leaks

#### Writing Files

```go
// Basic file writing
file, err := os.Create("output.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

// Write data
data := []byte("Hello, World!")
_, err = file.Write(data)
if err != nil {
    log.Fatal(err)
}

// Buffered writing
writer := bufio.NewWriter(file)
_, err = writer.WriteString("Hello, World!")
if err != nil {
    log.Fatal(err)
}
writer.Flush()
```

**Writing Considerations:**

1. **Atomicity**: Ensure write operations complete
2. **Buffering**: Balance memory and performance
3. **Permissions**: Set appropriate file modes
4. **Error Handling**: Handle write failures

## Working with Paths

### Understanding Path Operations

Go's `path/filepath` package provides tools for working with file paths in a cross-platform manner.

#### Path Manipulation

```go
// Join paths
path := filepath.Join("dir", "subdir", "file.txt")

// Get directory
dir := filepath.Dir(path)

// Get filename
filename := filepath.Base(path)

// Get extension
ext := filepath.Ext(path)

// Clean path
cleanPath := filepath.Clean("dir/../dir/./file.txt")
```

**Path Concepts:**

1. **Cross-Platform**: Works on all operating systems
2. **Normalization**: Handles path separators
3. **Validation**: Ensures path correctness
4. **Manipulation**: Build and modify paths

#### File Information

```go
// Get file info
info, err := os.Stat("file.txt")
if err != nil {
    log.Fatal(err)
}

// Check file type
isDir := info.IsDir()
size := info.Size()
modTime := info.ModTime()
mode := info.Mode()
```

**File Info Uses:**

1. **Validation**: Check file existence
2. **Metadata**: Access file properties
3. **Permissions**: Check access rights
4. **Timestamps**: Track modifications

## File Manipulation

### Understanding File Manipulation

Go provides tools for more complex file operations through the `os` package.

#### File Operations

```go
// Rename file
err := os.Rename("old.txt", "new.txt")

// Remove file
err := os.Remove("file.txt")

// Create directory
err := os.Mkdir("dir", 0755)

// Create all directories
err := os.MkdirAll("dir/subdir", 0755)
```

**Manipulation Concepts:**

1. **Atomicity**: Operations complete or fail
2. **Permissions**: Set appropriate modes
3. **Recursion**: Handle nested structures
4. **Error Handling**: Handle operation failures

#### File Permissions

```go
// Set permissions
err := os.Chmod("file.txt", 0644)

// Change owner
err := os.Chown("file.txt", uid, gid)

// Check permissions
file, err := os.OpenFile("file.txt", os.O_RDWR, 0644)
```

**Permission Concepts:**

1. **Unix-style**: Octal permission notation
2. **Security**: Appropriate access control
3. **Portability**: Cross-platform considerations
4. **Default Values**: Common permission patterns

## Advanced File Operations

### Understanding Advanced Operations

Go provides advanced file operations for complex scenarios.

#### File Locking

```go
type FileLocker struct {
    mu sync.Mutex
    f  *os.File
}

func (l *FileLocker) Lock() error {
    l.mu.Lock()
    return syscall.Flock(int(l.f.Fd()), syscall.LOCK_EX)
}

func (l *FileLocker) Unlock() error {
    defer l.mu.Unlock()
    return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}
```

**Locking Concepts:**

1. **Exclusive Access**: Prevent concurrent writes
2. **Deadlock Prevention**: Proper lock management
3. **Cross-Process**: Coordinate between processes
4. **Error Handling**: Handle lock failures

#### Memory-Mapped Files

```go
func mapFile(filename string) ([]byte, error) {
    f, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    fi, err := f.Stat()
    if err != nil {
        return nil, err
    }

    data, err := syscall.Mmap(int(f.Fd()), 0, int(fi.Size()),
        syscall.PROT_READ, syscall.MAP_SHARED)
    if err != nil {
        return nil, err
    }

    return data, nil
}
```

**Memory Mapping Concepts:**

1. **Performance**: Direct memory access
2. **Large Files**: Handle big datasets
3. **Shared Memory**: Process communication
4. **Resource Management**: Proper cleanup

## Best Practices

### File Operations

1. **Always Close Files**

   - Use `defer` for cleanup
   - Handle close errors
   - Prevent resource leaks

2. **Check Errors**

   - Validate all operations
   - Provide context
   - Handle failures gracefully

3. **Use Buffering**

   - Balance memory usage
   - Improve performance
   - Handle large files

4. **Set Permissions**
   - Follow security principles
   - Use appropriate modes
   - Consider portability

### Path Operations

1. **Use `filepath`**

   - Cross-platform compatibility
   - Path normalization
   - Safe manipulation

2. **Validate Paths**

   - Check existence
   - Verify permissions
   - Handle special cases

3. **Clean Paths**
   - Remove redundancies
   - Handle separators
   - Prevent traversal

### Error Handling

1. **Contextual Errors**

   - Include operation details
   - Provide recovery options
   - Log appropriately

2. **Resource Management**
   - Proper cleanup
   - Handle panics
   - Monitor resources

## Common Patterns

### File Processing Pipeline

```go
type FileProcessor struct {
    input  string
    output string
    stages []func([]byte) []byte
}

func (p *FileProcessor) Process() error {
    // Read input
    data, err := os.ReadFile(p.input)
    if err != nil {
        return err
    }

    // Apply stages
    for _, stage := range p.stages {
        data = stage(data)
    }

    // Write output
    return os.WriteFile(p.output, data, 0644)
}
```

### Concurrent File Processing

```go
type FileWorker struct {
    files  chan string
    result chan error
}

func (w *FileWorker) Process() {
    for file := range w.files {
        err := processFile(file)
        w.result <- err
    }
}

func processFiles(files []string, workers int) error {
    worker := &FileWorker{
        files:  make(chan string),
        result: make(chan error),
    }

    // Start workers
    for i := 0; i < workers; i++ {
        go worker.Process()
    }

    // Send files
    go func() {
        for _, file := range files {
            worker.files <- file
        }
        close(worker.files)
    }()

    // Collect results
    for range files {
        if err := <-worker.result; err != nil {
            return err
        }
    }
    return nil
}
```

Remember:

- Always handle errors properly
- Clean up resources
- Use appropriate buffering
- Consider performance implications
- Follow security best practices
- Test edge cases
- Document complex operations
- Consider concurrent access
- Use appropriate file modes
- Validate file operations

This guide covers the fundamental and advanced aspects of file operations in Go. Understanding these concepts is crucial for building robust and efficient applications.
