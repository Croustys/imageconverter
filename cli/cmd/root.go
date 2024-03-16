package webptopng

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/image/webp"
)

var version = "1.0.0"
var rootCmd = &cobra.Command{
	Use:     "webptopng",
	Version: version,
	Short:   "webptopng - Converts webp to png",
	Long:    `webptopng - Converts webp to png`,
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		fileName := filepath.Base(path)

		f0, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f0.Close()
		img0, err := webp.Decode(f0)
		if err != nil {
			fmt.Println(err)
			return
		}
		pngFile, err := os.Create("./" + fileNameWithoutExtension(fileName) + ".png")
		if err != nil {
			fmt.Println(err)
		}
		err = png.Encode(pngFile, img0)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
