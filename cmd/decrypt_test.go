package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecryptCmd_WithInMemoryInput(t *testing.T) {
	input := []byte("hello encrypted world")

	// Encrypt first to get a valid input stream
	var encOut bytes.Buffer
	enc := *encryptCmd
	enc.SetIn(bytes.NewReader(input))
	enc.SetOut(&encOut)
	enc.Flags().Set("in", "-")
	enc.Flags().Set("out", "-")
	enc.Flags().Set("password", "testpass")
	require.NoError(t, enc.RunE(&enc, []string{}))

	// Decrypt
	dec := *decryptCmd
	dec.SetIn(bytes.NewReader(encOut.Bytes()))
	var decOut bytes.Buffer
	dec.SetOut(&decOut)
	dec.Flags().Set("in", "-")
	dec.Flags().Set("out", "-")
	dec.Flags().Set("password", "testpass")

	err := dec.RunE(&dec, []string{})
	require.NoError(t, err)
	require.Equal(t, input, decOut.Bytes())
}

func TestDecryptCmd_WithFileInputOutput(t *testing.T) {
	// Create input file with encrypted data
	input := []byte("file roundtrip test")
	tmpDir := t.TempDir()

	// Write encrypted file
	inPath := tmpDir + "/encrypted.aes"
	outPath := tmpDir + "/decrypted.txt"

	encFile, err := os.Create(inPath)
	require.NoError(t, err)

	enc := *encryptCmd
	enc.SetIn(bytes.NewReader(input))
	enc.SetOut(encFile)
	enc.Flags().Set("in", "-")
	enc.Flags().Set("out", "-")
	enc.Flags().Set("password", "filepass")
	require.NoError(t, enc.RunE(&enc, []string{}))
	_ = encFile.Close()

	// Run decrypt
	dec := *decryptCmd
	dec.Flags().Set("in", inPath)
	dec.Flags().Set("out", outPath)
	dec.Flags().Set("password", "filepass")

	err = dec.RunE(&dec, []string{})
	require.NoError(t, err)

	// Validate output
	plain, err := os.ReadFile(outPath)
	require.NoError(t, err)
	require.Equal(t, input, plain)
}
