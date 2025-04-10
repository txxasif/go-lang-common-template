# Go Reflection: A Comprehensive Guide

## Table of Contents

1. [Introduction](#introduction)
2. [Type Reflection](#type-reflection)
3. [Value Reflection](#value-reflection)
4. [Struct Tags](#struct-tags)
5. [Basic Use Cases](#basic-use-cases)
6. [Advanced Reflection](#advanced-reflection)
7. [Best Practices](#best-practices)
8. [Common Patterns](#common-patterns)

## Introduction

Reflection in Go is a powerful feature that allows programs to examine and manipulate their own structure at runtime. While powerful, it should be used judiciously due to its performance implications and complexity.

### Why Reflection Matters

1. **Dynamic Programming**: Handle unknown types at runtime
2. **Serialization**: Convert between formats
3. **Validation**: Check data structures
4. **Code Generation**: Create code dynamically
5. **Framework Development**: Build flexible systems

## Type Reflection

### Understanding Type Reflection

Type reflection in Go allows you to examine and work with types at runtime through the `reflect` package.

#### Basic Type Inspection

```go
func inspectType(v interface{}) {
    t := reflect.TypeOf(v)

    // Type information
    fmt.Printf("Type: %v\n", t)
    fmt.Printf("Kind: %v\n", t.Kind())
    fmt.Printf("Name: %v\n", t.Name())
    fmt.Printf("Package: %v\n", t.PkgPath())
}
```

**Key Concepts:**

1. **Type Objects**: Represent Go types
2. **Kind**: Basic type category
3. **Name**: Type name
4. **Package**: Defining package

#### Advanced Type Operations

```go
func typeOperations(t reflect.Type) {
    // Check if type implements interface
    var r io.Reader
    readerType := reflect.TypeOf(&r).Elem()
    fmt.Printf("Implements Reader: %v\n", t.Implements(readerType))

    // Get method information
    for i := 0; i < t.NumMethod(); i++ {
        method := t.Method(i)
        fmt.Printf("Method: %v\n", method.Name)
    }

    // Get field information
    if t.Kind() == reflect.Struct {
        for i := 0; i < t.NumField(); i++ {
            field := t.Field(i)
            fmt.Printf("Field: %v\n", field.Name)
        }
    }
}
```

**Advanced Concepts:**

1. **Interface Implementation**: Check interface satisfaction
2. **Method Inspection**: Examine type methods
3. **Field Analysis**: Study struct fields
4. **Type Comparison**: Compare types

## Value Reflection

### Understanding Value Reflection

Value reflection allows you to examine and manipulate values at runtime.

#### Basic Value Operations

```go
func inspectValue(v interface{}) {
    val := reflect.ValueOf(v)

    // Value information
    fmt.Printf("Value: %v\n", val)
    fmt.Printf("Kind: %v\n", val.Kind())
    fmt.Printf("Type: %v\n", val.Type())
    fmt.Printf("CanSet: %v\n", val.CanSet())
}
```

**Key Concepts:**

1. **Value Objects**: Represent Go values
2. **Kind**: Value's type category
3. **Type**: Value's type
4. **Settability**: Can value be modified

#### Advanced Value Manipulation

```go
func manipulateValue(v interface{}) {
    val := reflect.ValueOf(v)

    // Create new value
    newVal := reflect.New(val.Type())

    // Copy value
    if val.CanInterface() {
        newVal.Elem().Set(val)
    }

    // Modify value
    if val.CanSet() {
        switch val.Kind() {
        case reflect.Int:
            val.SetInt(42)
        case reflect.String:
            val.SetString("hello")
        }
    }
}
```

**Advanced Concepts:**

1. **Value Creation**: Create new values
2. **Value Copying**: Duplicate values
3. **Value Modification**: Change values
4. **Type Conversion**: Convert between types

## Struct Tags

### Understanding Struct Tags

Struct tags are string literals that provide metadata about struct fields.

#### Basic Tag Usage

```go
type User struct {
    ID        int    `json:"id" db:"user_id"`
    Name      string `json:"name" db:"user_name"`
    Email     string `json:"email" db:"user_email"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

**Key Concepts:**

1. **Tag Syntax**: Key-value pairs
2. **Tag Parsing**: Extract metadata
3. **Common Uses**: Serialization, validation
4. **Tag Conventions**: Standard formats

#### Advanced Tag Operations

```go
func processTags(v interface{}) {
    t := reflect.TypeOf(v)
    if t.Kind() != reflect.Struct {
        return
    }

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)

        // Get tag value
        jsonTag := field.Tag.Get("json")
        dbTag := field.Tag.Get("db")

        // Parse tag options
        jsonName, jsonOpts := parseTag(jsonTag)
        dbName, dbOpts := parseTag(dbTag)

        // Process tags
        processField(field, jsonName, jsonOpts, dbName, dbOpts)
    }
}
```

**Advanced Concepts:**

1. **Tag Parsing**: Extract and validate
2. **Tag Options**: Handle additional metadata
3. **Tag Validation**: Check tag correctness
4. **Tag Processing**: Apply tag rules

## Basic Use Cases

### Understanding Common Applications

Reflection is commonly used in several standard scenarios.

#### Serialization

```go
func marshalJSON(v interface{}) ([]byte, error) {
    val := reflect.ValueOf(v)
    if val.Kind() != reflect.Struct {
        return nil, fmt.Errorf("expected struct")
    }

    result := make(map[string]interface{})
    t := val.Type()

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        jsonTag := field.Tag.Get("json")
        if jsonTag == "" {
            continue
        }

        fieldVal := val.Field(i)
        result[jsonTag] = fieldVal.Interface()
    }

    return json.Marshal(result)
}
```

**Serialization Concepts:**

1. **Field Mapping**: Map struct to format
2. **Tag Processing**: Handle metadata
3. **Value Conversion**: Convert types
4. **Error Handling**: Validate data

#### Validation

```go
func validateStruct(v interface{}) error {
    val := reflect.ValueOf(v)
    if val.Kind() != reflect.Struct {
        return fmt.Errorf("expected struct")
    }

    t := val.Type()
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        fieldVal := val.Field(i)

        // Check required
        if required := field.Tag.Get("required"); required == "true" {
            if isZero(fieldVal) {
                return fmt.Errorf("%s is required", field.Name)
            }
        }

        // Validate field
        if err := validateField(field, fieldVal); err != nil {
            return err
        }
    }

    return nil
}
```

**Validation Concepts:**

1. **Field Validation**: Check field values
2. **Tag Rules**: Apply validation rules
3. **Zero Values**: Detect empty fields
4. **Error Reporting**: Provide feedback

## Advanced Reflection

### Understanding Advanced Techniques

Go's reflection system supports complex operations for advanced use cases.

#### Dynamic Function Calls

```go
func callFunction(fn interface{}, args ...interface{}) ([]interface{}, error) {
    v := reflect.ValueOf(fn)
    if v.Kind() != reflect.Func {
        return nil, fmt.Errorf("not a function")
    }

    // Convert arguments
    in := make([]reflect.Value, len(args))
    for i, arg := range args {
        in[i] = reflect.ValueOf(arg)
    }

    // Call function
    out := v.Call(in)

    // Convert results
    results := make([]interface{}, len(out))
    for i, val := range out {
        results[i] = val.Interface()
    }

    return results, nil
}
```

**Advanced Concepts:**

1. **Function Reflection**: Examine functions
2. **Dynamic Calls**: Call functions dynamically
3. **Argument Handling**: Process arguments
4. **Result Processing**: Handle returns

#### Type Creation

```go
func createType(kind reflect.Kind, fields []reflect.StructField) reflect.Type {
    // Create struct type
    structType := reflect.StructOf(fields)

    // Create slice type
    if kind == reflect.Slice {
        return reflect.SliceOf(structType)
    }

    // Create pointer type
    if kind == reflect.Ptr {
        return reflect.PtrTo(structType)
    }

    return structType
}
```

**Creation Concepts:**

1. **Type Construction**: Build types
2. **Field Definition**: Define struct fields
3. **Type Composition**: Combine types
4. **Type Validation**: Ensure correctness

## Best Practices

### Reflection Usage

1. **Performance Considerations**

   - Avoid in hot paths
   - Cache reflection results
   - Use type assertions when possible
   - Profile before optimizing

2. **Type Safety**

   - Validate types early
   - Handle errors properly
   - Use type switches
   - Document assumptions

3. **Code Organization**
   - Isolate reflection code
   - Use helper functions
   - Document thoroughly
   - Test extensively

### Error Prevention

1. **Type Errors**

   - Check types early
   - Handle edge cases
   - Validate inputs
   - Provide clear errors

2. **Performance Issues**
   - Monitor usage
   - Cache results
   - Use alternatives
   - Profile regularly

## Common Patterns

### Dynamic Struct Creation

```go
func createStruct(fields map[string]interface{}) interface{} {
    var structFields []reflect.StructField

    for name, value := range fields {
        structFields = append(structFields, reflect.StructField{
            Name: strings.Title(name),
            Type: reflect.TypeOf(value),
        })
    }

    structType := reflect.StructOf(structFields)
    return reflect.New(structType).Interface()
}
```

### Tag Processing

```go
func processTags(tag string) (name string, options map[string]string) {
    parts := strings.Split(tag, ",")
    name = parts[0]
    options = make(map[string]string)

    for _, part := range parts[1:] {
        if strings.Contains(part, "=") {
            kv := strings.Split(part, "=")
            options[kv[0]] = kv[1]
        } else {
            options[part] = ""
        }
    }

    return name, options
}
```

Remember:

- Use reflection sparingly
- Understand performance impact
- Handle errors properly
- Document thoroughly
- Test extensively
- Cache when possible
- Validate inputs
- Consider alternatives
- Profile regularly
- Follow best practices

This guide covers the fundamental and advanced aspects of Go's reflection system. Understanding these concepts is crucial for building flexible and dynamic Go applications.
