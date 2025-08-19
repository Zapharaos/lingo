package test

import (
	"os"
	"testing"
)

// checkErrDeleteFolder Check if an error is not nil, delete the folder and fail the test
func checkErrDeleteFolder(t *testing.T, err error, dname string) {
	if err == nil {
		return
	}
	t.Error(err)
	err = os.RemoveAll(dname)
	if err != nil {
		t.Error(err)
	}
	t.FailNow()
}

type Suite struct {
	dname   string // temporary directory name (and new working directory during the test)
	oldWd   string // old working directory (to go back to it after the test)
	created bool   // whether the temporary directory has been created
}

func NewSuite() *Suite {
	return &Suite{
		dname:   "",
		oldWd:   "",
		created: false,
	}
}

// Create a full test suite
// Please clean test suite after use (defer ts.Clean())
func (s *Suite) Create(t *testing.T) (directoryPath string) {
	// Create temporary directory
	dname, err := os.MkdirTemp("", "lingo-test-suite-")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	s.created = true
	s.dname = dname

	// Save the old working directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	s.oldWd = oldWd

	// Chdir to a temporary directory
	err = os.Chdir(dname)
	checkErrDeleteFolder(t, err, dname)

	// Setup config directory
	s.createConfigFiles(t)

	return dname
}

func (s *Suite) createConfigFiles(t *testing.T) {
	// Create config translations directory
	err := os.MkdirAll("config/translations", os.ModePerm)
	checkErrDeleteFolder(t, err, s.dname)

	// Create the default english translation file
	file, err := os.Create("config/translations/active.en.toml")
	checkErrDeleteFolder(t, err, s.dname)

	// Write the translation content
	_, err = file.WriteString(`
		[hello]
		other = "Hello, {{.name}}!"
	`)
	checkErrDeleteFolder(t, err, s.dname)

	// Close the file
	err = file.Close()
	checkErrDeleteFolder(t, err, s.dname)

	// French translation file
	file, err = os.Create("config/translations/active.fr.toml")
	checkErrDeleteFolder(t, err, s.dname)
	err = file.Close()
	checkErrDeleteFolder(t, err, s.dname)
}

// Clean a test suite
func (s *Suite) Clean(t *testing.T) {
	// go back to the old working directory
	err := os.Chdir(s.oldWd)
	if err != nil {
		t.Error(err)
	}

	// remove the temporary directory
	err = os.RemoveAll(s.dname)
	if err != nil {
		t.Error(err)
	}
}
