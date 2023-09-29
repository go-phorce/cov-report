package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatters_Keys(t *testing.T) {
	k := keys(formatters)
	exp := []string{"ds", "json", "txt", "xml"}
	assert.Equal(t, exp, k)
}

func Test_ColorOf(t *testing.T) {
	assert.Equal(t, colorOf(0), "red")
	assert.Equal(t, colorOf(79), "red")
	assert.Equal(t, colorOf(81), "peru")
	assert.Equal(t, colorOf(89), "peru")
	assert.Equal(t, colorOf(90), "green")
	assert.Equal(t, colorOf(99), "green")
}

var testResult = &coverageResult{
	coverageAccumulator: coverageAccumulator{
		Total:           1234,
		Covered:         1053,
		Uncovered:       181,
		Excluded:        5,
		ExcludedSources: []string{"codegenc.go"},
	},
	Result:          85.33,
	ResultFormatted: `85.33%`,
	TopUncovered: []fileResult{
		{Filename: "bob.go", Total: 100, Covered: 10, Uncovered: 90, CoveredPct: 10.0, UncoveredPct: 90.0},
	},
}

func loadResult(t *testing.T) (string, string) {
	f := filepath.Join("testdata", strings.ToLower(t.Name())) + ".result"
	res, err := os.ReadFile(f)
	require.NoError(t, err)
	return string(res), f
}

func verifyFormatter(t *testing.T, fName string, verifier func(exp, act, msg string)) {
	b := bytes.Buffer{}
	formatters[fName](&b, testResult)
	res, fn := loadResult(t)
	verifier(res, b.String(), fn)
}

func TestFormatters_JSON(t *testing.T) {
	verifyFormatter(t, "json", func(exp, act, msg string) { assert.JSONEq(t, exp, act, msg) })
}

func TestFormatters_XML(t *testing.T) {
	verifyFormatter(t, "xml", func(exp, act, msg string) { assert.Equal(t, exp, act, msg) })
}

func TestFormatters_DS(t *testing.T) {
	verifyFormatter(t, "ds", func(exp, act, msg string) { assert.Equal(t, exp, act, msg) })
}

func TestFormatters_Txt(t *testing.T) {
	verifyFormatter(t, "txt", func(exp, act, msg string) { assert.Equal(t, exp, act, msg) })
}
