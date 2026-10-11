module github.com/abh/certman

go 1.26.0

require (
	github.com/fsnotify/fsnotify v1.10.1
	github.com/pkg/errors v0.9.1
)

require golang.org/x/sys v0.49.0 // indirect

// v0.3.0 was published with the wrong module path (github.com/dyson/certmon).
retract v0.3.0
