package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileResult_Finish(t *testing.T) {
	r := fileResult{Total: 1000, Covered: 980}
	r.finish()
	assert.EqualValues(t, 98.0, r.CoveredPct)

	r = fileResult{}
	r.finish()
	assert.EqualValues(t, 0.0, r.CoveredPct)
}

func TestFileResult_Longest(t *testing.T) {
	f := []fileResult{
		{Filename: "bob.go"},
		{Filename: "alice.go"},
		{Filename: "eve.go"},
	}
	assert.Equal(t, 8, longestFilename(f))
	assert.Equal(t, 0, longestFilename(nil))
	assert.Equal(t, 0, longestFilename([]fileResult{}))
}
