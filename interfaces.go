package autoupdate

// AutoUpdate defines methods to be filled by a special structure that may be
// used for Lantern to update itself.
type AutoUpdate interface {
	// SetVersion sets the internal release number of the current process.
	SetVersion(int)

	// Version returns the internal release number of the current process.
	Version() int

	// Query sends the current software checksum to an update server. If the
	// update server decides this program is outdated, it will send information
	// on how to update. With this information, a Patch can be constructed.
	Query() (Patch, error)

	// Do will periodically check for updates (using Query()) without
	// interrupting the main process. In an update is found it will download and
	// apply it without user interaction.
	Do()
}

// Patch defines methods for applying binary patches.
type Patch interface {
	// Apply downloads and apply the binary diff against the actual program file.
	// The returned value will be nil if, and only if, we're absolutely sure the
	// update was applied successfully.
	Apply() error

	// Version returns the internal release number of the update.
	Version() int
}
