package version

import "fmt"

var (
	Version   string
	BuildDate string
	BuildRef  string
)

func VersionString() string {
	return fmt.Sprintf("Tapper version %s buildref %s date %s",
		Version, BuildRef, BuildDate,
	)
}

func PrintVersion() {
	fmt.Println(VersionString())
}
