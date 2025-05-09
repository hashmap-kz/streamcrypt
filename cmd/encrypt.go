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

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Compress and encrypt an input stream",
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath, _ := cmd.Flags().GetString("in")
		outputPath, _ := cmd.Flags().GetString("out")
		password, _ := cmd.Flags().GetString("password")

		// Use stdin if inputPath is not set
		var in io.Reader
		if inputPath == "" || inputPath == "-" {
			in = cmd.InOrStdin()
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
			out = cmd.OutOrStdout()
		} else {
			f, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer f.Close()
			out = f
		}

		// Compression setup
		compressor := codec.GzipCompressor{}

		// Encryption setup
		crypter := aesgcm.NewChunkedGCMCrypter(password)

		r, err := pipe.CompressAndEncryptOptional(in, compressor, crypter)
		if err != nil {
			return fmt.Errorf("pipeline setup failed: %w", err)
		}

		_, err = io.Copy(out, r)
		if err != nil {
			return fmt.Errorf("error during encryption: %w", err)
		}

		return nil
	},
}

func init() {
	encryptCmd.Flags().StringP("in", "i", "", "Input file (default: stdin)")
	encryptCmd.Flags().StringP("out", "o", "", "Output file (default: stdout)")
	encryptCmd.Flags().StringP("password", "p", "", "Password to derive encryption key")
	_ = encryptCmd.MarkFlagRequired("password")
}
