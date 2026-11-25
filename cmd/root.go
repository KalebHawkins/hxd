/*
Copyright © 2026 KalebHawkins@outlook.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const (
	bufferSize = 16
)

var (
	filename          string
	startAddress      string
	numberBytesToDump int
)

var rootCmd = &cobra.Command{
	Use:   "hxd -f <file> [-s offset] [-n length]",
	Short: "Dump files as hex with offsets and printable ASCII",
	Long: `hxd prints the contents of a file as a hex dump: one line per 16 bytes,
showing the offset, the bytes in hex, and the printable ASCII characters
(non-printable bytes appear as '.').

Use --offset to start reading at a specific position and --length to limit
how many bytes are dumped. Offsets accept decimal or prefixed values
(0x40, 0o100, 0b1000000), so addresses can be passed as written.`,
	Example: `  hxd -f game.exe                 dump the whole file
  hxd -f game.exe -n 64           first 64 bytes (headers)
  hxd -f game.exe -s 0x40 -n 64   64 bytes starting at offset 0x40
  hxd -f save.dat > save.txt      save a dump to diff against later`,
	SilenceUsage: true,
	Args:         cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		if numberBytesToDump < 0 {
			return fmt.Errorf("-n must not be negative")
		}

		var v int64
		var err error
		if startAddress != "" {
			v, err = strconv.ParseInt(startAddress, 0, 64)
			if err != nil {
				return fmt.Errorf("invalid --offset [-s]: %q, %w", startAddress, err)
			}
		}

		if v < 0 {
			return fmt.Errorf("invalid --offset [-s], offset cannot be negative: %q", startAddress)
		}

		f, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer f.Close()

		ret, err := f.Seek(v, io.SeekStart)
		if err != nil {
			return fmt.Errorf("failed to seek: %w", err)
		}

		var r io.Reader = f
		if numberBytesToDump > 0 {
			r = io.LimitReader(r, int64(numberBytesToDump))
		}

		if err := dump(cmd.OutOrStdout(), r, ret); err != nil {
			return fmt.Errorf("failed to dump file: %w", err)
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&filename, "file", "f", "", "The file you want to dump")
	rootCmd.MarkFlagRequired("file")

	rootCmd.Flags().IntVarP(&numberBytesToDump, "length", "n", 0, "The number of bytes to dump (0 = all)")
	rootCmd.Flags().StringVarP(&startAddress, "offset", "s", "", "The offset to start at.")
}

func writeHex(sb *strings.Builder, buf []byte, numberOfBytes int) {
	for i := 0; i < bufferSize; i++ {
		if i < numberOfBytes {
			fmt.Fprintf(sb, "%02x ", buf[i])
		} else {
			sb.WriteString("   ") // Three spaces to fill empty space.
		}
		if (i+1)%4 == 0 {
			sb.WriteString("| ")
		}
	}
}

func writeAscii(sb *strings.Builder, buf []byte, numberOfBytes int) {
	for _, b := range buf[:numberOfBytes] {
		// ASCII Character Range 0x20 -> 0x7e
		if b < 0x20 || b > 0x7e {
			sb.WriteString(".")
			continue
		}
		sb.WriteByte(b)
	}
}

func dump(w io.Writer, r io.Reader, startOffset int64) error {
	offset := startOffset
	buf := make([]byte, bufferSize)
	var sb strings.Builder

	for {
		n, err := r.Read(buf)

		if n > 0 {
			fmt.Fprintf(&sb, "0x%08x: ", offset)
			writeHex(&sb, buf, n)
			writeAscii(&sb, buf, n)
			sb.WriteByte('\n')
			if _, err := io.WriteString(w, sb.String()); err != nil {
				return err
			}
			sb.Reset()
		}
		offset += int64(n)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

	}

	return nil
}
