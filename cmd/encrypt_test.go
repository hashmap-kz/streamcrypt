package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptCmd_WithInMemoryInput(t *testing.T) {
	input := []byte("hello world")

	cmd := *encryptCmd
	cmd.SetIn(bytes.NewReader(input)) // instead of os.Stdin
	var out bytes.Buffer
	cmd.SetOut(&out)

	cmd.Flags().Set("in", "-")
	cmd.Flags().Set("out", "-")
	cmd.Flags().Set("password", "secret")

	err := cmd.RunE(&cmd, []string{})
	require.NoError(t, err)
	require.Greater(t, out.Len(), 0)
}

func TestEncryptCmd_WithFileInputOutput(t *testing.T) {
	inputData := []byte("file-based input encryption test")

	// Create temp input file
	inFile, err := os.CreateTemp(t.TempDir(), "input-*.txt")
	require.NoError(t, err)
	defer inFile.Close()
	_, err = inFile.Write(inputData)
	require.NoError(t, err)

	// Create temp output file path (will be created by command)
	outFilePath := inFile.Name() + ".enc"

	// Clone and configure the command
	cmd := *encryptCmd
	cmd.Flags().Set("in", inFile.Name())
	cmd.Flags().Set("out", outFilePath)
	cmd.Flags().Set("password", "testpass")
	err = cmd.RunE(&cmd, []string{})
	require.NoError(t, err)

	// Check the output file exists and is not empty
	outStat, err := os.Stat(outFilePath)
	require.NoError(t, err)
	require.Greater(t, outStat.Size(), int64(0), "output file should not be empty")

	// Clean up
	_ = os.Remove(outFilePath)
}
