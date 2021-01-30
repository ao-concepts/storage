# ao-concepts storage module

![CI](https://github.com/ao-concepts/storage/workflows/CI/badge.svg)
[![codecov](https://codecov.io/gh/ao-concepts/storage/branch/master/graph/badge.svg?token=AQVUZTRGQS)](https://codecov.io/gh/ao-concepts/storage)

This module provides some useful helpers for efficient work with [gorm.io/gorm](gorm.io/gorm).

## Information

The ao-concepts ecosystem is still under active development and therefore the API of this module may have breaking changes until there is a first stable release.

If you are interested in contributing to this project, feel free to open a issue to discus a new feature, enhancement or improvement. If you found a bug or security vulnerability in this package, please start a issue, or open a PR against `master`.

## Installation

```
go get -u github.com/ao-concepts/storage
```

## Storage

The storage `Controller` is the storage provider for your application.
It is used to create transactions that access yout storage system.
The log parameter is optional. You cann pass a custom logger there.

You can use the `Controller` to start transactions.
There is also a `UseTransaction` function that can be used to wrap a
function that should be executed within a transaction. 
The transaction will be rolled back if the error returned is not nil.

```go
dialector := sqlite.Open(":memory:")

c := storage.New(dialector, nil)
```

## Repositories

The `Repository` struct is intended to be embedded into a custom repository struct:

```go
import (
   "github.com/ao-concepts/storage"
   "gorm.io/gorm"
)

type User struct {
   gorm.Model
   Name string
}

type UserRepository struct {
   storage.Repository
}

func (r *UserRepository) GetAll(tx *storage.Transaction) (users []User, err error) {
   return tx.Gorm().Find(&users).Error
}
```

## Used packages 

This project uses some really great packages. Please make sure to check them out!

| Package                                                                              | Usage                             |
| ------------------------------------------------------------------------------------ | --------------------------------- |
| [github.com/stretchr/testify](https://github.com/stretchr/testify)                   | Testing                           |
| [gorm.io/gorm](gorm.io/gorm)                                                         | Database access                   |
