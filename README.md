# gord - Generic Repository Pattern with GORM

## Overview

`gord` is a generic repository pattern implementation in Go, built on top of GORM, a popular ORM library. It provides a clean and reusable interface for common CRUD (Create, Read, Update, Delete) operations, enabling a more structured and maintainable data access layer in your Go applications.

## Features

- **Generic CRUD Operations**: Provides a set of generic CRUD methods that work with any data model.
- **Type Safety**: Uses Go’s generics to enforce type safety for models and IDs.
- **GORM Integration**: Seamlessly integrates with GORM, leveraging its powerful features.
- **Custom Updates**: Allows for flexible updates with type-safe validation.
- **Bulk Operations**: Supports bulk delete and save operations.
- **Existence Check**: Provides a method to check if a record exists by its ID.

## Installation

First, ensure that you have GORM installed in your Go project:

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite # Replace with your database driver
```

Then, include this package in your project:

```bash
go get github.com/sanginsu/gord
```

## Usage

### 1. Define Your Models

Start by defining your models that you want to manage using the repository.

```go
package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Name  string
    Email string
}
```

### 2. Initialize the Repository

Create an instance of the repository for your model.

```go
package main

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/yourusername/gord"
    "yourapp/models"
)

func main() {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    db.AutoMigrate(&models.User{})

    userRepository := gord.Repository[models.User, uint]{tx: db}

    // Use userRepository for CRUD operations
}
```

### 3. Perform CRUD Operations

You can now use the repository instance to perform various operations.

```go
// Create a new user
user := models.User{Name: "John Doe", Email: "john.doe@example.com"}
err := userRepository.Create(user)

// Find a user by ID
foundUser, err := userRepository.FindByID(1)

// Update a user
updates := gord.UpdateMap{
    "Name": "Jane Doe",
}
err := userRepository.Updates(foundUser, updates)

// Delete a user
err := userRepository.Delete(foundUser)

// Check if a user exists
exists, err := userRepository.ExistByID(1)
```

### 4. Customizing the Repository

If you need custom methods, you can embed the `Repository` struct and add your methods.

```go
type UserRepository struct {
    gord.Repository[models.User, uint]
}

func (r UserRepository) FindByEmail(email string) (models.User, error) {
    var user models.User
    err := r.tx.Where("email = ?", email).First(&user).Error
    return user, err
}
```

### 5. Validating Updates

`gord` includes a validation mechanism to ensure that only allowed types are used in update operations.

```go
updates := gord.UpdateMap{
    "Name": "New Name",
    // "InvalidField": []int{1, 2, 3}, // This would cause a validation error
}

if err := updates.valid(); err != nil {
    fmt.Println("Validation failed:", err)
} else {
    userRepository.Updates(user, updates)
}
```

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss your ideas.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.
