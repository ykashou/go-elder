module github.com/ykashou/go-elder

go 1.24.2

require (
	// github.com/audiomage-dev/go-benchmarks v0.0.0
	github.com/spf13/cobra v1.8.0
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace (
	// Map canonical paths to local submodule directories in internal/
	// github.com/audiomage-dev/go-benchmarks => ./internal/go-benchmarks
	// Map canonical paths to local submodule directories in pkg/
	// github.com/audiomage-dev/go-magefile => ./pkg/go-magefile

// NOTE: May need replaces for transitive dependencies if they are also managed as submodules and referenced directly or indirectly within go-elder context.
)
