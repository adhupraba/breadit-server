Live application - https://breadit.adhupraba.com/

Frontend repo - https://github.com/adhupraba/breadit-client/

---

# Create a goose migration file

```bash
goose create users sql
```

# Migration

use the migrate.sh to migrate the schema to the database

```bash
sh migrate.sh
```

# SQLC generate

```bash
sqlc generate
```

# Running development server

This automatically restarts dev server on code changes

```bash
compiledaemon --command="./breadit-server"
```

# Build for production

```bash
go build -tags netgo -ldflags '-s -w' -o breadit-server
```

1. **`go build`**:
   This is the Go command to compile Go programs. By default, it compiles the Go code in the current directory (unless a specific package or file is mentioned).

2. **`-tags netgo`**:
   Build tags are a way to include/exclude certain files from the build process based on conditions.

   - `netgo`: This specific tag tells the Go compiler to use the Go-based `net` package instead of the system's native networking libraries. This can be helpful in situations where you want your application to be purely Go-based without relying on system C libraries for networking, potentially increasing portability and reducing issues with library dependencies.

3. **`-ldflags '-s -w'`**:
   The `-ldflags` (linker flags) option allows you to pass flags to the Go linker.

   - `-s`: Omit the symbol table and debug information.
   - `-w`: Omit the DWARF symbol table.

   Using these flags (`-s` and `-w`) reduces the size of the resulting binary by removing symbol tables and debug information. It's a common practice to use these flags when building release versions of applications where you want to minimize the binary size and aren't concerned with having debug information embedded in the binary.

4. **`-o app`**:
   The `-o` flag specifies the output file name for the compiled binary. In this case, the resulting binary will be named `app`.

In summary, this command compiles the Go code in the current directory using the Go-based `net` package, omits symbol and debug information to reduce binary size, and outputs the resulting binary with the name `app`.
