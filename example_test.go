package version_test

import (
	"fmt"

	"github.com/ardnew/version"
)

func init() {
	// list the version history directly in version.ChangeLog. the last element
	// is used as the current version of this package.
	version.ChangeLog = []version.Change{
		{
			Version: "0.1.0",
			Date:    "Feb 26, 2020", // very many date-time formats recognized
			Description: []string{
				`initial commit`,
			},
		}, {
			Version: "0.1.0+fqt",
			Title:   "Formal Test",
			Description: []string{
				`update user manual`,
			},
		}, {
			Version: "0.2.0-beta+red",
			Title:   "Red Label",
			Date:    "20-Mar-9 17:45:23",
			Description: []string{
				`add feature: Dude`,
				`fix bug: Sweet`,
			},
		},
	}
}

func Example() {

	// show that we are currently using the last entry in ChangeLog
	fmt.Printf("using ChangeLog version %q\n\n", version.String())

	// print a pretty changelog to stdout.
	version.PrintChanges()

	version.Set("0.1.4")
	fmt.Printf("set version to %q\n\n", version.String())

	// Output:
	// using ChangeLog version "0.2.0-beta+red"
	//
	// ――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	//  version 0.1.0                                    Wed, 26 Feb 2020 00:00:00 UTC
	// ――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	//   initial commit
	//
	// ――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	//  version 0.1.0+fqt - Formal Test
	// ――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	//   update user manual
	//
	// ――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	//  version 0.2.0-beta+red - Red Label               Mon, 09 Mar 2020 17:45:23 UTC
	// ――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――――
	//   add feature: Dude
	//   fix bug: Sweet
	//
	// set version to "0.1.4"
}
