[![Go Report Card](https://goreportcard.com/badge/github.com/nccapo/kv-store)](https://goreportcard.com/report/github.com/nccapo/kv-store)
![Go Version](https://img.shields.io/github/go-mod/go-version/nccapo/kv-store)
[![Go Reference](https://pkg.go.dev/badge/github.com/nccapo/kv-store.svg)](https://pkg.go.dev/github.com/nccapo/kv-store)
![License](https://img.shields.io/github/license/nccapo/kv-store)
![Issues Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)
![Version](https://img.shields.io/github/v/release/nccapo/kv-store)

# Key-Value Storage (Redis Clone) in Go 

## Overview

This project is an in-memory key-value storage system inspired by Redis , implemented in Go . The goal of this project is to provide a simple yet powerful tool for storing and retrieving data with features like expiration, atomic operations, and concurrency support. It is designed primarily as a learning exercise to gain more experience with Go programming and concurrent systems.

## Features
- **Thread-Safe** : Uses sync.Map for thread-safe operations.
- **Key Expiration** : Automatically expires keys after a specified duration.
- **Batch Operations** : Perform multiple Get or Set operations in a single call.
- **Persistence** : Save and load data from disk using snapshots.
- **Monitoring** : Collect metrics for operations like `Set`, `Get`, and `Delete`.

## Contributing
We welcome contributions from the community! Here's how you can help:
1. **Fork the Repository** : Click the "Fork" button on GitHub.
2. **Clone Your Fork** 
   ```shell
    git clone https://github.com/nccapo/kv-store.git
    ```
3. **Create a Branch** :
   ```shell
    git checkout -b feature/new-feature
    ```
4. **Submit a Pull Request** : Push your branch and open a PR against the main branch.

## License
This project is licensed under the MIT License. See the [LICENSE](https://github.com/nccapo/kv-store/blob/master/LICENSE) file for details.