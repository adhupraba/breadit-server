# Running development server

This automatically restarts dev server on code changes

```bash
compiledaemon --command="./breadit-server"
```

# Build for production

```bash
go build -tags netgo -ldflags '-s -w' -o breadit-server
```
