// Copyright (c) 2015 Jesse Meek <https://github.com/waigani>
// This program is Free Software see LICENSE file for details.

package diffparser_test

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"testing"

	"github.com/International/diffparser"
)

func rawDiff() string {
	byt, err := ioutil.ReadFile("example.diff")
	if err != nil {
		log.Fatal("could not find test file")
	}
	return string(byt)
}

// TODO(waigani) tests are missing more creative names (spaces, special
// chars), and diffed files that are not in the current directory.

func TestFileModeAndNaming(t *testing.T) {
	diff, err := diffparser.Parse(rawDiff())
	require.Nil(t, err)
	require.Len(t, diff.Files, 6)

	for i, expected := range []struct {
		mode     diffparser.FileMode
		origName string
		newName  string
	}{
		{
			mode:     diffparser.MODIFIED,
			origName: "file1",
			newName:  "file1",
		},
		{
			mode:     diffparser.DELETED,
			origName: "file2",
			newName:  "",
		},
		{
			mode:     diffparser.DELETED,
			origName: "file3",
			newName:  "",
		},
		{
			mode:     diffparser.NEW,
			origName: "",
			newName:  "file4",
		},
		{
			mode:     diffparser.NEW,
			origName: "",
			newName:  "newname",
		},
		{
			mode:     diffparser.DELETED,
			origName: "symlink",
			newName:  "",
		},
	} {
		file := diff.Files[i]
		t.Logf("testing file: %v", file)
		require.Equal(t, file.Mode, expected.mode)
		require.Equal(t, file.OrigName, expected.origName)
		require.Equal(t, file.NewName, expected.newName)
	}
}

func TestHunk(t *testing.T) {
	diff, err := diffparser.Parse(rawDiff())
	require.Nil(t, err)
	require.Len(t, diff.Files, 6)

	expectedOrigLines := []diffparser.DiffLine{
		{
			Mode:     diffparser.UNCHANGED,
			Number:   1,
			Content:  "some",
			Position: 2,
		}, {
			Mode:     diffparser.UNCHANGED,
			Number:   2,
			Content:  "lines",
			Position: 3,
		}, {
			Mode:     diffparser.REMOVED,
			Number:   3,
			Content:  "in",
			Position: 4,
		}, {
			Mode:     diffparser.UNCHANGED,
			Number:   4,
			Content:  "file1",
			Position: 5,
		},
	}

	expectedNewLines := []diffparser.DiffLine{
		{
			Mode:     diffparser.ADDED,
			Number:   1,
			Content:  "add a line",
			Position: 1,
		}, {
			Mode:     diffparser.UNCHANGED,
			Number:   2,
			Content:  "some",
			Position: 2,
		}, {
			Mode:     diffparser.UNCHANGED,
			Number:   3,
			Content:  "lines",
			Position: 3,
		}, {
			Mode:     diffparser.UNCHANGED,
			Number:   4,
			Content:  "file1",
			Position: 5,
		},
	}

	file := diff.Files[0]
	origRange := file.Hunks[0].OrigRange
	newRange := file.Hunks[0].NewRange

	require.Equal(t, origRange.Start, 1)
	require.Equal(t, origRange.Length, 4)
	require.Equal(t, newRange.Start, 1)
	require.Equal(t, newRange.Length, 4)

	for i, line := range expectedOrigLines {
		require.Equal(t, *origRange.Lines[i], line)
	}
	for i, line := range expectedNewLines {
		require.Equal(t, *newRange.Lines[i], line)
	}
}
