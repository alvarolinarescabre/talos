// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package block_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/siderolabs/talos/pkg/machinery/config/configloader"
	"github.com/siderolabs/talos/pkg/machinery/config/encoder"
	"github.com/siderolabs/talos/pkg/machinery/config/types/block"
	"github.com/siderolabs/talos/pkg/machinery/constants"
	blockres "github.com/siderolabs/talos/pkg/machinery/resources/block"
)

func TestUserVolumeConfigMarshalUnmarshal(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name string

		filename string
		cfg      func(t *testing.T) *block.UserVolumeConfigV1Alpha1
	}{
		{
			name:     "disk selector",
			filename: "uservolumeconfig_diskselector.yaml",
			cfg: func(t *testing.T) *block.UserVolumeConfigV1Alpha1 {
				c := block.NewUserVolumeConfigV1Alpha1()
				c.MetaName = "ceph-data"

				require.NoError(t, c.ProvisioningSpec.DiskSelectorSpec.Match.UnmarshalText([]byte(`disk.transport == "nvme" && !system_disk`)))
				c.ProvisioningSpec.ProvisioningMinSize = block.MustByteSize("10GiB")
				c.ProvisioningSpec.ProvisioningMaxSize = block.MustByteSize("100GiB")
				c.FilesystemSpec.FilesystemType = blockres.FilesystemTypeXFS

				return c
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cfg := test.cfg(t)

			warnings, err := cfg.Validate(validationMode{})
			require.NoError(t, err)
			require.Empty(t, warnings)

			marshaled, err := encoder.NewEncoder(cfg, encoder.WithComments(encoder.CommentsDisabled)).Encode()
			require.NoError(t, err)

			t.Log(string(marshaled))

			expectedMarshaled, err := os.ReadFile(filepath.Join("testdata", test.filename))
			require.NoError(t, err)

			assert.Equal(t, string(expectedMarshaled), string(marshaled))

			provider, err := configloader.NewFromBytes(expectedMarshaled)
			require.NoError(t, err)

			docs := provider.Documents()
			require.Len(t, docs, 1)

			assert.Equal(t, cfg, docs[0])
		})
	}
}

func TestUserVolumeConfigValidate(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name string

		cfg func(t *testing.T) *block.UserVolumeConfigV1Alpha1

		expectedErrors string
	}{
		{
			name: "no name",

			cfg: func(t *testing.T) *block.UserVolumeConfigV1Alpha1 {
				c := block.NewUserVolumeConfigV1Alpha1()

				require.NoError(t, c.ProvisioningSpec.DiskSelectorSpec.Match.UnmarshalText([]byte(`disk.size > 1u`)))
				c.ProvisioningSpec.ProvisioningMinSize = block.MustByteSize("2.5TiB")

				return c
			},

			expectedErrors: "name is required\nname must be between 1 and 34 characters long",
		},
		{
			name: "too long name",

			cfg: func(t *testing.T) *block.UserVolumeConfigV1Alpha1 {
				c := block.NewUserVolumeConfigV1Alpha1()
				c.MetaName = strings.Repeat("X", 35)

				require.NoError(t, c.ProvisioningSpec.DiskSelectorSpec.Match.UnmarshalText([]byte(`disk.size > 1u`)))
				c.ProvisioningSpec.ProvisioningMinSize = block.MustByteSize("2.5TiB")

				return c
			},

			expectedErrors: "name must be between 1 and 34 characters long",
		},
		{
			name: "invalid characters in name",

			cfg: func(t *testing.T) *block.UserVolumeConfigV1Alpha1 {
				c := block.NewUserVolumeConfigV1Alpha1()
				c.MetaName = "some/name"

				require.NoError(t, c.ProvisioningSpec.DiskSelectorSpec.Match.UnmarshalText([]byte(`disk.size > 1u`)))
				c.ProvisioningSpec.ProvisioningMinSize = block.MustByteSize("2.5TiB")

				return c
			},

			expectedErrors: "name can only contain lowercase and uppercase ASCII letters, digits, and hyphens",
		},
		{
			name: "invalid disk selector",

			cfg: func(t *testing.T) *block.UserVolumeConfigV1Alpha1 {
				c := block.NewUserVolumeConfigV1Alpha1()
				c.MetaName = constants.EphemeralPartitionLabel

				require.NoError(t, c.ProvisioningSpec.DiskSelectorSpec.Match.UnmarshalText([]byte(`disk.size > 120`)))

				return c
			},

			expectedErrors: "disk selector is invalid: ERROR: <input>:1:11: found no matching overload for '_>_' applied to '(uint, int)'\n | disk.size > 120\n | ..........^\nmin size or max size is required",
		},
		{
			name: "min size greater than max size",

			cfg: func(t *testing.T) *block.UserVolumeConfigV1Alpha1 {
				c := block.NewUserVolumeConfigV1Alpha1()
				c.MetaName = constants.EphemeralPartitionLabel

				c.ProvisioningSpec.ProvisioningMinSize = block.MustByteSize("2.5TiB")
				c.ProvisioningSpec.ProvisioningMaxSize = block.MustByteSize("10GiB")

				return c
			},

			expectedErrors: "disk selector is required\nmin size is greater than max size",
		},
		{
			name: "unsupported filesystem type",

			cfg: func(t *testing.T) *block.UserVolumeConfigV1Alpha1 {
				c := block.NewUserVolumeConfigV1Alpha1()
				c.MetaName = constants.EphemeralPartitionLabel

				require.NoError(t, c.ProvisioningSpec.DiskSelectorSpec.Match.UnmarshalText([]byte(`disk.size > 120u * GiB`)))
				c.ProvisioningSpec.ProvisioningMaxSize = block.MustByteSize("2.5TiB")
				c.ProvisioningSpec.ProvisioningMinSize = block.MustByteSize("10GiB")
				c.FilesystemSpec.FilesystemType = blockres.FilesystemTypeISO9660

				return c
			},

			expectedErrors: "unsupported filesystem type: iso9660",
		},
		{
			name: "valid",

			cfg: func(t *testing.T) *block.UserVolumeConfigV1Alpha1 {
				c := block.NewUserVolumeConfigV1Alpha1()
				c.MetaName = constants.EphemeralPartitionLabel

				require.NoError(t, c.ProvisioningSpec.DiskSelectorSpec.Match.UnmarshalText([]byte(`disk.size > 120u * GiB`)))
				c.ProvisioningSpec.ProvisioningMaxSize = block.MustByteSize("2.5TiB")
				c.ProvisioningSpec.ProvisioningMinSize = block.MustByteSize("10GiB")
				c.FilesystemSpec.FilesystemType = blockres.FilesystemTypeEXT4

				return c
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			cfg := test.cfg(t)

			_, err := cfg.Validate(validationMode{})

			if test.expectedErrors == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)

				assert.EqualError(t, err, test.expectedErrors)
			}
		})
	}
}
