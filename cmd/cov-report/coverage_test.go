package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Min(t *testing.T) {
	assert.Equal(t, min(1, 1), 1)
	assert.Equal(t, min(10, 123), 10)
	assert.Equal(t, min(123, 10), 10)
}

func Test_Max(t *testing.T) {
	assert.Equal(t, max(1, 1), 1)
	assert.Equal(t, max(123, 1), 123)
	assert.Equal(t, max(123, 1234), 1234)
}

func Test_BlockKeys(t *testing.T) {
	b := blockKeys{
		{filename: "bob.go", startLine: 5, startCol: 6, endLine: 9, endCol: 10},
		{filename: "bob.go", startLine: 10, startCol: 1, endLine: 20, endCol: 10},
		{filename: "bob.go", startLine: 5, startCol: 1, endLine: 5, endCol: 5},
		{filename: "alice.go", startLine: 5, startCol: 1, endLine: 9, endCol: 10},
	}
	sort.Sort(b)
	exp := blockKeys{
		{filename: "alice.go", startLine: 5, startCol: 1, endLine: 9, endCol: 10},
		{filename: "bob.go", startLine: 5, startCol: 1, endLine: 5, endCol: 5},
		{filename: "bob.go", startLine: 5, startCol: 6, endLine: 9, endCol: 10},
		{filename: "bob.go", startLine: 10, startCol: 1, endLine: 20, endCol: 10},
	}
	assert.Equal(t, exp, b)
}

func Test_Main(t *testing.T) {
	verifyMain(t, 0, "-fmt", "txt", "testdata/cp1.out", "testdata/cp2.out", "testdata/cp3.out")
}

func Test_MainNoFiles(t *testing.T) {
	verifyMain(t, 2, "-fmt", "txt")
}

func Test_MainBadFormat(t *testing.T) {
	verifyMain(t, 2, "-fmt", "bob", "testdata/cp1.out")
}

func Test_MainMissingFile(t *testing.T) {
	verifyMain(t, 1, "testdata/cp1.out", "testdata/cp_missing.out")
}

func Test_MainBadFlags(t *testing.T) {
	verifyMain(t, 2, "-foo", "bar")
}

func verifyMain(t *testing.T, expResultCode int, args ...string) {
	fullArgs := append([]string{"cov-report"}, args...)
	out := bytes.Buffer{}
	res := realMain(&nopWriteCloser{&out}, fullArgs)
	assert.Equal(t, expResultCode, res)

	exp, fn := loadResult(t)
	str := out.String()
	assert.Equal(t, strings.Replace(exp, "\t", "    ", -1), strings.Replace(str, "\t", "    ", -1), fn)
}

type nopWriteCloser struct {
	w io.Writer
}

func (n *nopWriteCloser) Write(d []byte) (int, error) {
	return n.w.Write(d)
}

func (n *nopWriteCloser) Close() error {
	return nil
}

func Test_CoverageAccumulator(t *testing.T) {
	a := newCoverageAccumulator()
	for _, f := range []string{"cp1.out", "cp2.out", "cp3.out"} {
		err := a.parse(filepath.Join("testdata", f), nil)
		assert.NoError(t, err)
	}
	a.calculate()

	tf, err := os.CreateTemp("", t.Name())
	require.NoError(t, err)
	defer os.Remove(tf.Name())

	a.dumpCombined(tf.Name())
	assertTextFileEqual(t, "testdata/test_coverageaccumulater_cp.result", tf.Name())

	r := a.result(1)
	exp, fn := loadResult(t)
	b := bytes.Buffer{}
	formatters["json"](&b, r)
	assert.JSONEq(t, exp, b.String(), fn)
}

func assertTextFileEqual(t *testing.T, exp, act string) {
	e, err := os.ReadFile(exp)
	require.NoError(t, err)
	a, err := os.ReadFile(act)
	require.NoError(t, err)
	elines := strings.Split(string(e), "\n")
	alines := strings.Split(string(a), "\n")
	assert.Equal(t, elines, alines)
}
