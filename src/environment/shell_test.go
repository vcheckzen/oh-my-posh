package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalHostName(t *testing.T) {
	hostName := "hello"
	assert.Equal(t, hostName, cleanHostName(hostName))
}

func TestHostNameWithLocal(t *testing.T) {
	hostName := "hello.local"
	assert.Equal(t, "hello", cleanHostName(hostName))
}

func TestHostNameWithLan(t *testing.T) {
	hostName := "hello.lan"
	cleanHostName := cleanHostName(hostName)
	assert.Equal(t, "hello", cleanHostName)
}

func TestDirMatchesOneOf(t *testing.T) {
	cases := []struct {
		GOOS     string
		HomeDir  string
		Dir      string
		Pattern  string
		Expected bool
	}{
		{GOOS: LINUX, HomeDir: "/home/bill", Dir: "/home/bill", Pattern: "/home/bill", Expected: true},
		{GOOS: LINUX, HomeDir: "/home/bill", Dir: "/home/bill/foo", Pattern: "~/foo", Expected: true},
		{GOOS: LINUX, HomeDir: "/home/bill", Dir: "/home/bill/foo", Pattern: "~/Foo", Expected: false},
		{GOOS: LINUX, HomeDir: "/home/bill", Dir: "/home/bill/foo", Pattern: "~\\\\foo", Expected: true},
		{GOOS: LINUX, HomeDir: "/home/bill", Dir: "/home/bill/foo/bar", Pattern: "~/fo.*", Expected: true},
		{GOOS: LINUX, HomeDir: "/home/bill", Dir: "/home/bill/foo", Pattern: "~/fo\\w", Expected: true},

		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill", Pattern: "C:\\\\Users\\\\Bill", Expected: true},
		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill", Pattern: "C:/Users/Bill", Expected: true},
		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill", Pattern: "c:/users/bill", Expected: true},
		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill", Pattern: "~", Expected: true},
		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill\\Foo", Pattern: "~/Foo", Expected: true},
		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill\\Foo", Pattern: "~/foo", Expected: true},
		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill\\Foo\\Bar", Pattern: "~/fo.*", Expected: true},
		{GOOS: WINDOWS, HomeDir: "C:\\Users\\Bill", Dir: "C:\\Users\\Bill\\Foo", Pattern: "~/fo\\w", Expected: true},
	}

	for _, tc := range cases {
		got := dirMatchesOneOf(tc.Dir, tc.HomeDir, tc.GOOS, []string{tc.Pattern})
		assert.Equal(t, tc.Expected, got)
	}
}

func TestDirMatchesOneOfRegexInverted(t *testing.T) {
	// detect panic(thrown by MustCompile)
	defer func() {
		if err := recover(); err != nil {
			// display a message explaining omp failed(with the err)
			assert.Equal(t, "regexp: Compile(`^(?!Projects[\\/]).*$`): error parsing regexp: invalid or unsupported Perl syntax: `(?!`", err)
		}
	}()
	_ = dirMatchesOneOf("Projects/oh-my-posh", "", LINUX, []string{"(?!Projects[\\/]).*"})
}

func TestDirMatchesOneOfRegexInvertedNonEscaped(t *testing.T) {
	// detect panic(thrown by MustCompile)
	defer func() {
		if err := recover(); err != nil {
			// display a message explaining omp failed(with the err)
			assert.Equal(t, "regexp: Compile(`^(?!Projects/).*$`): error parsing regexp: invalid or unsupported Perl syntax: `(?!`", err)
		}
	}()
	_ = dirMatchesOneOf("Projects/oh-my-posh", "", LINUX, []string{"(?!Projects/).*"})
}
