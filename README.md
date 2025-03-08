# StoneKVR 🪨🚀

**StoneKVR** is an **embedded NoSQL key/value store** written in Go. It provides on‑disk persistence with simple query capabilities, making it a lightweight and efficient solution for modern applications.

## Overview 🔍

StoneKVR offers:

- **Embedded storage** for quick local data persistence.
- **On‑disk durability** ensuring your data stays safe.
- **Simple query capabilities** for effortless data retrieval.
- **Native Go integration** for fast and concurrent operations.

## Features ✨

- **Lightweight & Fast**: Optimized for performance in resource‑constrained environments.
- **Simple API**: Easy-to‑use methods to store, retrieve, and query data.
- **On‑Disk Persistence**: Reliable storage without the need for a separate server.
- **Concurrent Access**: Built with Go’s concurrency in mind.

## Quick Start 🚀

1. **Install via Go Modules:**

   ```bash
   go get github.com/cryptrunner49/stonekvr
   ```  

2. **Basic Usage Example:**

   ```go
   package main  

   import (
       "fmt"
       "log"

       "github.com/cryptrunner49/stonekvr"
   )

   func main() {
       // Initialize the store (creates a new one if it doesn't exist)
       stone, err := stonekvr.NewStone("data.db")
       if err != nil {
           log.Fatal(err)
       }
       defer stone.Close()

       // Set a key/value pair
       err = store.Set("greeting", "Hello, StoneKVR! 👋")
       if err != nil {
           log.Fatal(err)
       }

       // Retrieve a value
       value, err := stone.Get("greeting")
       if err != nil {
           log.Fatal(err)
       }
       fmt.Println(value)
   }
   ```

## Contributing 🤝

We welcome contributions! Check out our [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

## License 📜

Distributed under the MIT License. See [LICENSE](./LICENSE) for more information.
