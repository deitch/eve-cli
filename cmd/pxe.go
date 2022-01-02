// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/lf-edge/eve-cli/pkg/github"
	"github.com/lf-edge/eve-cli/pkg/pxe"
	"github.com/spf13/cobra"
)

const (
	org            = "lf-edge"
	repo           = "eve"
	defaultArch    = "amd64"
	defaultVersion = "latest"
	defaultBaseURL = "https://github.com/lf-edge/eve/releases/download"
)

var (
	outpath, serial, generateSerial, url, version, arch string
)

var pxeCmd = &cobra.Command{
	Use:   "pxe",
	Short: "generate iPXE boot script",
	Long:  `Generate iPXE boot script. Can be configured to point to any download location, including local. Also can generate soft serial.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if url != "" && ((arch != "" && cmd.Flags().Changed("arch")) || (version != "" && cmd.Flags().Changed("version"))) {
			return errors.New("cannot provide --arch or --version when --url is provided")
		}
		if serial != "" && generateSerial != "" && cmd.Flags().Changed("generate-serial") {
			return errors.New("cannot provide both --serial and --generate-serial")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		// figure out where to write the output
		var out io.Writer = os.Stdout
		if outpath != "" {
			outfile, err := os.OpenFile(outpath, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				log.Fatalf("failed to open file %s: %v", outpath, err)
			}
			defer outfile.Close()
			out = outfile
		}
		// generate a serial, if required
		if serial == "" {
			switch generateSerial {
			case "mac":
				serial = "${mac:hexhyp}"
			case "ip":
				serial = "${ip}"
			case "uuid":
				serial = uuid.New().String()
			default:
				log.Fatalf("unknown serial type '%s'", generateSerial)
			}
		}
		// figure out the base URL
		if url == "" {
			if version == "" || version == defaultVersion{
				version, err = github.GetLatestVersion(org, repo)
				if err != nil {
					log.Fatal(err)
				}
			}
			url = fmt.Sprintf("%s/%s/%s", defaultBaseURL, version, arch)
		}

		// generate the iPXE script
		content, err := pxe.GeneratePXE(serial, url)
		if err != nil {
			log.Fatalf("%v", err)
		}
		// write the output
		if _, err := out.Write(content); err != nil {
			log.Fatal(err)
		}
	},
}

func pxeInit() {
	pxeCmd.Flags().StringVar(&outpath, "out", "", "path where to store the iPXE script, blank to stdout")
	pxeCmd.Flags().StringVar(&serial, "serial", "", "provided serial to use")
	pxeCmd.Flags().StringVar(&generateSerial, "generate-serial", "uuid", "serial type to generate, one of: ip, mac, uuid; must provide either --serial or --generate-serial")
	pxeCmd.Flags().StringVar(&url, "url", "", fmt.Sprintf("base URL for assets, such as initrd, installer, rootfs. Defaults to latest release from %s for architecture %s", defaultBaseURL, defaultArch))
	pxeCmd.Flags().StringVar(&version, "version", defaultVersion, fmt.Sprintf("which version to take from the default base for assets on GitHub, when no URL is provided. '%s' means discovering the latest version. Conflicts with 'url'", defaultVersion))
	pxeCmd.Flags().StringVar(&arch, "arch", defaultArch, "which architecture to take from the default base for assets on GitHub, when no URL is provided. Conflicts with 'url'")
}
