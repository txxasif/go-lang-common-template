# OOP: Go vs JavaScript - A Comparative Guide

This guide compares Object-Oriented Programming (OOP) in **Go** and **JavaScript**, highlighting how each language approaches core OOP concepts like encapsulation, composition, and polymorphism. Go emphasizes simplicity and explicitness, while JavaScript offers flexibility with dynamic typing. Understanding these differences helps you write idiomatic code in either language.

## Table of Contents

1. [Introduction](#introduction)
2. [Type Definition](#type-definition)
3. [Constructors](#constructors)
4. [Methods](#methods)
5. [Inheritance vs Composition](#inheritance-vs-composition)
6. [Interfaces and Abstraction](#interfaces-and-abstraction)
7. [Encapsulation](#encapsulation)
8. [Properties](#properties)
9. [Error Handling](#error-handling)
10. [Polymorphism](#polymorphism)
11. [Testing OOP Code](#testing-oop-code)
12. [Best Practices](#best-practices)

## Introduction

**Go** and **JavaScript** approach OOP differently:

- **Go**: Uses structs, composition, and implicit interfaces for simplicity and compile-time safety. It avoids classical inheritance, favoring explicitness.
- **JavaScript**: Leverages classes, prototypes, and dynamic typing for flexibility, supporting both inheritance and composition.

This guide compares their approaches to help you choose the right patterns for your projects.

## Type Definition

### JavaScript: Classes

JavaScript uses the `class` keyword (ES6+) for object blueprints, built on prototype-based inheritance.

**Example:**

```javascript
class User {
  constructor(firstName, lastName) {
    this.firstName = firstName;
    this.lastName = lastName;
  }
}
```

### Go: Structs

Go uses `struct` to define data structures, with no concept of classes.

**Example:**

```go
package main

type User struct {
    FirstName string
    LastName  string
}
```

**Key Differences:**

- JavaScript classes are dynamic and prototype-based.
- Go structs are static, requiring explicit type declarations.
- JavaScript includes constructors in class definitions; Go uses separate functions.

## Constructors

### JavaScript

Constructors initialize instances within a `class`.

**Example:**

```javascript
class User {
  constructor(firstName, lastName, email) {
    this.firstName = firstName;
    this.lastName = lastName;
    this.email = email;
  }

  static fromJSON(data) {
    return new User(data.firstName, data.lastName, data.email);
  }
}

const user = User.fromJSON({
  firstName: "Alice",
  lastName: "Smith",
  email: "alice@example.com",
});
```

### Go

Go uses convention-based constructor functions, typically named `New<Type>`.

**Example:**

```go
package main

import "fmt"

type User struct {
    FirstName string
    LastName  string
    Email     string
}

func NewUser(firstName, lastName, email string) *User {
    return &User{
        FirstName: firstName,
        LastName:  lastName,
        Email:     email,
    }
}

func main() {
    user := NewUser("Alice", "Smith", "alice@example.com")
    fmt.Printf("%+v\n", user)
}
```

**Key Differences:**

- JavaScript constructors are part of the class; Go uses standalone functions.
- Go typically returns pointers (`*User`) for mutability.
- JavaScript supports factory methods via `static`; Go uses regular functions.

## Methods

### JavaScript

Methods are defined inside a class and implicitly bind to `this`.

**Example:**

```javascript
class User {
  constructor(firstName, lastName) {
    this.firstName = firstName;
    this.lastName = lastName;
  }

  getFullName() {
    return `${this.firstName} ${this.lastName}`;
  }

  setEmail(email) {
    this.email = email;
  }
}
```

### Go

Methods are functions with a receiver, declared outside the struct.

**Example:**

```go
package main

import "fmt"

type User struct {
    FirstName string
    LastName  string
    Email     string
}

func (u User) GetFullName() string {
    return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func (u *User) SetEmail(email string) {
    u.Email = email
}

func main() {
    user := &User{FirstName: "Bob", LastName: "Jones"}
    fmt.Println(user.GetFullName()) // Bob Jones
    user.SetEmail("bob@example.com")
}
```

**Key Differences:**

- Go requires explicit receivers (`u User` or `u *User`); JavaScript uses `this`.
- Go distinguishes value vs. pointer receivers for immutability vs. mutation.
- JavaScript methods are scoped to the class; Go methods are package-scoped.

## Inheritance vs Composition

### JavaScript: Inheritance

JavaScript supports classical inheritance via `extends`.

**Example:**

```javascript
class Entity {
  constructor() {
    this.id = Math.random();
    this.createdAt = new Date();
  }
}

class User extends Entity {
  constructor(firstName, lastName) {
    super();
    this.firstName = firstName;
    this.lastName = lastName;
  }
}
```

### Go: Composition

Go uses **embedding** for composition, not inheritance.

**Example:**

```go
package main

import (
    "fmt"
    "time"
)

type Entity struct {
    ID        int
    CreatedAt time.Time
}

type User struct {
    Entity
    FirstName string
    LastName  string
}

func main() {
    user := User{
        Entity:    Entity{ID: 1, CreatedAt: time.Now()},
        FirstName: "Charlie",
        LastName:  "Brown",
    }
    fmt.Printf("ID: %d, Name: %s\n", user.ID, user.FirstName) // Promoted field
}
```

**Key Differences:**

- JavaScript uses `extends` and `super` for inheritance; Go embeds structs.
- Go composition avoids tight coupling of inheritance.
- JavaScript supports method overriding; Go requires explicit method definitions.

## Interfaces and Abstraction

### JavaScript: Abstraction

JavaScript mimics abstract classes using errors or base classes.

**Example:**

```javascript
class Repository {
  create() {
    throw new Error("Must implement create");
  }
}

class UserRepository extends Repository {
  create(user) {
    console.log(`Saving user: ${user.firstName}`);
    return user;
  }
}
```

### Go: Interfaces

Go uses implicit interfaces for abstraction.

**Example:**

```go
package main

import (
    "errors"
    "fmt"
)

type Repository interface {
    Create(user any) error
}

type UserRepository struct{}

func (r UserRepository) Create(user any) error {
    u, ok := user.(User)
    if !ok {
        return errors.New("invalid user type")
    }
    fmt.Printf("Saving user: %s\n", u.FirstName)
    return nil
}

type User struct {
    FirstName string
}

func main() {
    var repo Repository = UserRepository{}
    user := User{FirstName: "Dave"}
    repo.Create(user)
}
```

**Key Differences:**

- Go interfaces are satisfied implicitly; JavaScript requires explicit extension.
- Go interfaces are lightweight and flexible; JavaScript abstract classes are heavier.
- Go avoids inheritance hierarchies; JavaScript embraces them.

## Encapsulation

### JavaScript: Encapsulation

JavaScript uses private fields (`#`) for encapsulation.

**Example:**

```javascript
class User {
  #email;

  constructor(email) {
    this.#email = email;
  }

  getEmail() {
    return this.#email;
  }
}
```

### Go: Encapsulation

Go uses capitalization for package-level access control.

**Example:**

```go
package main

import "fmt"

type User struct {
    Email     string // Exported
    password  string // Unexported
}

func (u *User) GetPassword() string {
    return u.password
}

func main() {
    user := &User{Email: "eve@example.com", password: "secret"}
    fmt.Println(user.Email)        // Accessible
    // fmt.Println(user.password)  // Compile error: unexported
    fmt.Println(user.GetPassword()) // Accessible via method
}
```

**Key Differences:**

- Go's encapsulation is package-based (capitalized = exported).
- JavaScript's `#` fields are instance-private.
- Go relies on conventions; JavaScript enforces privacy syntactically.

## Properties

### JavaScript

JavaScript supports getters and setters for controlled access.

**Example:**

```javascript
class User {
  #name;

  set name(value) {
    this.#name = value.trim();
  }

  get name() {
    return this.#name || "Unknown";
  }
}
```

### Go

Go uses explicit methods for getter/setter behavior.

**Example:**

```go
package main

import (
    "fmt"
    "strings"
)

type User struct {
    name string
}

func (u *User) SetName(name string) {
    u.name = strings.TrimSpace(name)
}

func (u User) Name() string {
    if u.name == "" {
        return "Unknown"
    }
    return u.name
}

func main() {
    user := &User{}
    user.SetName("  Frank  ")
    fmt.Println(user.Name()) // Frank
}
```

**Key Differences:**

- JavaScript has native getter/setter syntax; Go uses methods.
- Go method names follow conventions (e.g., `Name`, `SetName`).
- JavaScript properties integrate with assignment; Go requires explicit calls.

## Error Handling

### JavaScript: Error Handling

JavaScript uses exceptions with `try/catch`.

**Example:**

```javascript
class UserService {
  async createUser(data) {
    try {
      if (!data.email) {
        throw new Error("Email is required");
      }
      // Simulate DB save
      return { id: 1, ...data };
    } catch (err) {
      throw new Error(`Failed to create user: ${err.message}`);
    }
  }
}
```

### Go: Error Handling

Go treats errors as values, returned explicitly.

**Example:**

```go
package main

import (
    "errors"
    "fmt"
)

type UserService struct{}

func (s UserService) CreateUser(data User) (User, error) {
    if data.Email == "" {
        return User{}, errors.New("email is required")
    }
    // Simulate DB save
    user := User{ID: 1, Email: data.Email}
    return user, nil
}

type User struct {
    ID    int
    Email string
}

func main() {
    svc := UserService{}
    user, err := svc.CreateUser(User{Email: ""})
    if err != nil {
        fmt.Println("Error:", err) // Error: email is required
        return
    }
    fmt.Printf("Created: %+v\n", user)
}
```

**Key Differences:**

- JavaScript uses exceptions; Go uses return values.
- Go requires explicit error checks; JavaScript centralizes with `try/catch`.
- Go errors are composable with `fmt.Errorf` and `%w`.

## Polymorphism

### JavaScript: Polymorphism

JavaScript achieves polymorphism through inheritance or interfaces (via classes).

**Example:**

```javascript
class Animal {
  speak() {
    return "Sound";
  }
}

class Dog extends Animal {
  speak() {
    return "Woof";
  }
}

const animals = [new Animal(), new Dog()];
animals.forEach((a) => console.log(a.speak())); // Sound, Woof
```

### Go: Polymorphism

Go uses interfaces for polymorphic behavior.

**Example:**

```go
package main

import "fmt"

type Animal interface {
    Speak() string
}

type Dog struct{}

func (Dog) Speak() string {
    return "Woof"
}

type Cat struct{}

func (Cat) Speak() string {
    return "Meow"
}

func main() {
    animals := []Animal{Dog{}, Cat{}}
    for _, a := range animals {
        fmt.Println(a.Speak()) // Woof, Meow
    }
}
```

**Key Differences:**

- Go interfaces are implicit and structural; JavaScript relies on class hierarchies.
- Go polymorphism is lightweight; JavaScript is more flexible but complex.
- Go avoids runtime type ambiguity; JavaScript embraces it.

## Testing OOP Code

### JavaScript: Testing

Use frameworks like Jest to test classes.

**Example:**

```javascript
// user.test.js
const { User } = require("./user");

test("User.getFullName concatenates names", () => {
  const user = new User("Alice", "Smith");
  expect(user.getFullName()).toBe("Alice Smith");
});

class User {
  constructor(firstName, lastName) {
    this.firstName = firstName;
    this.lastName = lastName;
  }

  getFullName() {
    return `${this.firstName} ${this.lastName}`;
  }
}
```

### Go: Testing

Go's built-in `testing` package simplifies unit tests.

**Example:**

```go
package main

import (
    "testing"
)

type User struct {
    FirstName string
    LastName  string
}

func (u User) GetFullName() string {
    return u.FirstName + " " + u.LastName
}

func TestUser_GetFullName(t *testing.T) {
    user := User{FirstName: "Alice", LastName: "Smith"}
    got := user.GetFullName()
    want := "Alice Smith"
    if got != want {
        t.Errorf("Got %q, want %q", got, want)
    }
}
```

**Key Differences:**

- Go tests are lightweight and built-in; JavaScript relies on external frameworks.
- Go emphasizes table-driven tests; JavaScript uses assertions like `expect`.
- Go's static typing reduces runtime test failures.

## Best Practices

### JavaScript: Best Practices

- **Use Dependency Injection:**
  ```javascript
  class UserService {
    constructor(repository) {
      this.repository = repository;
    }
  }
  ```
- **Leverage Method Chaining:**
  ```javascript
  class QueryBuilder {
    #condition;
    where(condition) {
      this.#condition = condition;
      return this;
    }
  }
  ```
- **Keep Classes Focused**: Single responsibility per class.
- **Use Private Fields**: Prefer `#field` for encapsulation.
- **Handle Errors Gracefully**: Centralize error handling with `try/catch`.

### Go: Best Practices

- **Explicit Dependencies:**

  ```go
  type UserService struct {
      repo Repository
  }

  func NewUserService(repo Repository) *UserService {
      return &UserService{repo: repo}
  }
  ```

- **Enable Method Chaining:**

  ```go
  type QueryBuilder struct {
      condition string
  }

  func (q *QueryBuilder) Where(condition string) *QueryBuilder {
      q.condition = condition
      return q
  }
  ```

- **Small Interfaces**: Define minimal interfaces for flexibility.
- **Use Pointer Receivers Judiciously**: For mutation or large structs.
- **Check Errors Explicitly**: Avoid ignoring return values.

**Key Differences:**

- Go favors explicitness (e.g., pointers, errors); JavaScript abstracts complexity.
- JavaScript encourages dynamic patterns; Go prioritizes compile-time safety.
- Go avoids deep hierarchies; JavaScript supports them.

## Key Takeaways

- **Go**:
  - Simplifies OOP with structs, interfaces, and composition.
  - Emphasizes explicitness, type safety, and minimalism.
  - Ideal for performance-critical, concurrent systems.
- **JavaScript**:
  - Offers flexible OOP with classes, inheritance, and dynamic typing.
  - Suits rapid prototyping and web development.
  - Trades safety for expressiveness.
- **Common Ground**:
  - Both support encapsulation and polymorphism.
  - Dependency injection and method chaining are universal.
- **Choosing Between Them**:
  - Use Go for backend systems needing reliability and scalability.
  - Use JavaScript for frontends or projects requiring flexibility.

This guide illustrates how Go and JavaScript achieve OOP goals differently, empowering you to pick the right approach for your needs.

_Last Updated: April 11, 2025_
