package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/hashmap-kz/streaming-compress-encrypt/pkg/crypt/aesgcm"

	"github.com/hashmap-kz/streaming-compress-encrypt/pkg/codec"
	"github.com/hashmap-kz/streaming-compress-encrypt/pkg/pipe"
	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt and decompress an input stream",
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath, _ := cmd.Flags().GetString("in")
		outputPath, _ := cmd.Flags().GetString("out")
		password, _ := cmd.Flags().GetString("password")

		// Use stdin if inputPath is not set
		var in io.Reader
		if inputPath == "" || inputPath == "-" {
			in = os.Stdin
		} else {
			f, err := os.Open(inputPath)
			if err != nil {
				return fmt.Errorf("failed to open input file: %w", err)
			}
			defer f.Close()
			in = f
		}

		// Use stdout if outputPath is not set
		var out io.Writer
		if outputPath == "" || outputPath == "-" {
			out = os.Stdout
		} else {
			f, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer f.Close()
			out = f
		}

		// Decompression setup
		decompressor := codec.GzipDecompressor{}

		// Decryption setup
		crypter := aesgcm.NewChunkedGCMCrypter(password)

		r, err := pipe.DecryptAndDecompressOptional(in, crypter, decompressor)
		if err != nil {
			return fmt.Errorf("pipeline setup failed: %w", err)
		}
		defer r.Close()

		_, err = io.Copy(out, r)
		if err != nil {
			return fmt.Errorf("error during decryption: %w", err)
		}

		return nil
	},
}

func init() {
	decryptCmd.Flags().StringP("in", "i", "", "Input file (default: stdin)")
	decryptCmd.Flags().StringP("out", "o", "", "Output file (default: stdout)")
	decryptCmd.Flags().StringP("password", "p", "", "Password to derive encryption key")
	_ = decryptCmd.MarkFlagRequired("password")
}
