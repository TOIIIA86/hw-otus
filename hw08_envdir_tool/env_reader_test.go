package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testDir            = "testdata/test_dir"
	testDirWithBadFile = "testdata/test_dir_bad"
	badDirPath         = "testdata/123123/"
	badFileName        = "123=123"

	successCaseDirPath = "testdata/env"
)

func TestReadDir(t *testing.T) {
	defer cleanUp()
	setUp()

	t.Run("Empty dir", func(t *testing.T) {
		res, err := ReadDir(testDir)
		require.NoError(t, err)
		require.Len(t, res, 0)
	})

	t.Run("Non existent dir", func(t *testing.T) {
		res, err := ReadDir(badDirPath)
		require.Error(t, err)
		require.Equal(t, true, os.IsNotExist(err))
		require.Len(t, res, 0)
	})

	t.Run("Bad file name err", func(t *testing.T) {
		_, err := os.Create(testDirWithBadFile + "/" + badFileName)
		if err != nil {
			log.Fatal(err)
		}

		res, err := ReadDir(testDirWithBadFile)

		require.Error(t, err)
		require.EqualError(t, err, ErrInvalidFileName.Error())
		require.Len(t, res, 0)
	})

	t.Run("Success case", func(t *testing.T) {
		expectedRes := Environment{
			"BAR": EnvVal{
				Value: "bar",
			},
			"EMPTY": EnvVal{
				UnsetVal: true,
			},
			"FOO": EnvVal{
				Value: "   foo\nwith new line",
			},
			"HELLO": EnvVal{
				Value: `"hello"`,
			},
			"UNSET": EnvVal{
				UnsetVal: true,
			},
		}

		res, err := ReadDir(successCaseDirPath)

		require.NoError(t, err)
		require.Equal(t, expectedRes, res)
	})
}

func cleanUp() {
	if err := os.RemoveAll(testDir); err != nil {
		log.Fatal(err)
	}

	if err := os.RemoveAll(testDirWithBadFile); err != nil {
		log.Fatal(err)
	}
}

func setUp() {
	err := os.Mkdir(testDir, 0o755)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir(testDirWithBadFile, 0o755)
	if err != nil {
		log.Fatal(err)
	}
}
