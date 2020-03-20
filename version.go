package version

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

// Version is the current version of the package. Use Set() or define ChangeLog
// to set the version.
var Version struct {
	Major      uint
	Minor      uint
	Patch      uint
	Prerelease string
	Metadata   string
}

// VersionPattern defines the regular expression used to validate and identify
// the components of a semantic version string.
//
// Source: https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
var VersionPattern = `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

// See `go doc time.Parse` for formatting convention.
var (
	// DateFormat defines the recognized date string formats, parsed in order of
	// of decreasing precedence.
	DateFormat = []string{
		`2006 Jan 2`,
		`2006-Jan-2`,
		`2006-1-2`,
		`2006 1 2`,
		`1-2-2006`,
		`1/2/2006`,

		`Jan 2, 2006`,

		`06 Jan 2`,
		`06-Jan-2`,
		`06-1-2`,
		`06 1 2`,
		`1-2-06`,
		`1/2/06`,

		`Jan 2, 06`,
	}

	// TimeFormat defines the recognized time string formats, parsed in order of
	// of decreasing precedence.
	TimeFormat = []string{
		`15:04:05`,

		`03:04:05PM`,
		`03:04:05pm`,
		`3:04:05PM`,
		`3:04:05pm`,

		`15:04`,

		`03:04PM`,
		`03:04pm`,
		`3:04PM`,
		`3:04pm`,
	}

	// DateTimeFormat defines the format used to write the date-time of a version
	// change; the output format string.
	DateTimeFormat = time.RFC1123
)

// Change represents the details of a version change.
type Change struct {
	Version     string
	Title       string
	Date        string
	Description []string
}

// ParseDate parses the given date-time string. It attempts every permutation of
// each DateFormat and TimeFormat pair (in either order), returning the first
// successfully-parsed time.Time object. If none of the pairs are successful,
// each DateFormat (ignoring TimeFormat) is then attempted. Finally, each of the
// standard formats provided by the time package are attempted.
func ParseDate(date string) *time.Time {
	if "" != date {
		for _, fd := range DateFormat {
			for _, ft := range TimeFormat {
				dt := fmt.Sprintf("%s %s", fd, ft)
				if t, err := time.Parse(dt, date); nil == err {
					return &t
				}
				td := fmt.Sprintf("%s %s", ft, fd)
				if t, err := time.Parse(td, date); nil == err {
					return &t
				}
			}
		}
		for _, fd := range DateFormat {
			if t, err := time.Parse(fd, date); nil == err {
				return &t
			}
		}
		for _, dt := range []string{
			time.ANSIC,       // "Mon Jan _2 15:04:05 2006"
			time.UnixDate,    // "Mon Jan _2 15:04:05 MST 2006"
			time.RubyDate,    // "Mon Jan 02 15:04:05 -0700 2006"
			time.RFC822,      // "02 Jan 06 15:04 MST"
			time.RFC822Z,     // "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
			time.RFC850,      // "Monday, 02-Jan-06 15:04:05 MST"
			time.RFC1123,     // "Mon, 02 Jan 2006 15:04:05 MST"
			time.RFC1123Z,    // "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
			time.RFC3339,     // "2006-01-02T15:04:05Z07:00"
			time.RFC3339Nano, // "2006-01-02T15:04:05.999999999Z07:00"
		} {
			if t, err := time.Parse(dt, date); nil == err {
				return &t
			}
		}
	}
	return nil
}

// String returns a formatted, multi-line string describing Change c.
func (c *Change) String() string {
	Parse(c.Version) // validate version string. will panic if invalid.

	const (
		maxWidth = 80
		titlePad = 1
		descPad  = 2
	)

	runeRepeat := func(c rune, n int) string {
		b := strings.Builder{}
		for i := 0; i < n; i++ {
			b.WriteRune(c)
		}
		return b.String()
	}

	// construct the "version - title" left-hand side
	vsb := strings.Builder{}
	vsb.WriteString("version ")
	vsb.WriteString(c.Version)
	if "" != c.Title {
		vsb.WriteString(" - ")
		vsb.WriteString(c.Title)
	}

	// construct the "date" right-hand side
	dsb := strings.Builder{}
	if t := ParseDate(c.Date); nil != t {
		dsb.WriteString(t.Format(DateTimeFormat))
	}

	// calculate the padding width between left- and right-hand sides
	middlePad := maxWidth - ((vsb.Len() + titlePad) + (dsb.Len() + titlePad))

	// horizontal line used for containing the header
	horizLine := runeRepeat('―', maxWidth) + "\n"

	// construct the header containing horizontal lines, version, title, and date
	b := strings.Builder{}
	b.WriteString(horizLine)
	b.WriteString(fmt.Sprintf("%*s%s", titlePad, "", vsb.String()))
	if dsb.Len() > 0 {
		b.WriteString(fmt.Sprintf("%*s%s", middlePad, "", dsb.String()))
	}
	b.WriteRune('\n')
	b.WriteString(horizLine)

	// append each description line with indentation
	if nil != c.Description {
		for _, line := range c.Description {
			fmt.Fprintf(&b, "%*s%s\n", descPad, "", line)
		}
	}

	return b.String()
}

// ChangeLog contains the history of version changes.
var ChangeLog []Change

// Parse validates a semantic version string and returns each of its components.
// It panics if the given version string is invalid.
func Parse(version string) (major, minor, patch uint, pre, meta string) {
	re := regexp.MustCompile(VersionPattern)
	sub := re.FindStringSubmatch(version)
	if 0 == len(sub) {
		panic("invalid version: " + version)
	}
	fmt.Sscanf(sub[1], "%d", &major)
	fmt.Sscanf(sub[2], "%d", &minor)
	fmt.Sscanf(sub[3], "%d", &patch)
	if len(sub) > 4 && "" != sub[4] {
		pre = sub[4]
	}
	if len(sub) > 5 && "" != sub[5] {
		meta = sub[5]
	}
	return
}

// Set sets the package version using a given semantic version string.
// It panics if the given version string is invalid.
func Set(version string) {
	Version.Major, Version.Minor, Version.Patch,
		Version.Prerelease, Version.Metadata = Parse(version)
}

// IsSet returns true if and only if the package version has been set.
// The package version is considered not-set if all components are equal to
// their zero value.
func IsSet() bool {
	return Version.Major != 0 || Version.Minor != 0 || Version.Patch != 0 ||
		Version.Prerelease != "" || Version.Metadata != ""
}

// String returns the semantic version string of the package.
// If the version has not been set, the last entry in ChangeLog is used (or
// panics if the last entry in ChangeLog contains an invalid version string).
// If ChangeLog has also not been set, an empty string is returned.
func String() string {

	str := func(major uint, minor uint, patch uint, pre string, meta string) string {
		b := strings.Builder{}
		fmt.Fprintf(&b, "%d.%d.%d", major, minor, patch)
		if "" != pre {
			b.WriteRune('-')
			b.WriteString(pre)
		}
		if "" != meta {
			b.WriteRune('+')
			b.WriteString(meta)
		}
		return b.String()
	}

	if IsSet() {
		return str(Version.Major, Version.Minor, Version.Patch,
			Version.Prerelease, Version.Metadata)
	} else if nil != ChangeLog && len(ChangeLog) > 0 {
		return str(Parse(ChangeLog[len(ChangeLog)-1].Version))
	}
	return ""
}

// FprintChangeLog writes to given io.Writer w all of the entries in ChangeLog.
// Panics if any of the entries have invalid version strings.
func FprintChangeLog(w io.Writer) {
	if nil != ChangeLog {
		for _, c := range ChangeLog {
			fmt.Fprintf(w, "%s\n", c.String())
		}
	}
}

// PrintChanges writes to stdout all of the entries in ChangeLog.
// Panics if any of the entries have invalid version strings.
func PrintChangeLog() {
	FprintChangeLog(os.Stdout)
}
