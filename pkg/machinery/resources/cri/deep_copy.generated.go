// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Code generated by "deep-copy -type RegistriesConfigSpec -type ImageCacheConfigSpec -type SeccompProfileSpec -header-file ../../../../hack/boilerplate.txt -o deep_copy.generated.go ."; DO NOT EDIT.

package cri

// DeepCopy generates a deep copy of RegistriesConfigSpec.
func (o RegistriesConfigSpec) DeepCopy() RegistriesConfigSpec {
	var cp RegistriesConfigSpec = o
	if o.RegistryMirrors != nil {
		cp.RegistryMirrors = make(map[string]*RegistryMirrorConfig, len(o.RegistryMirrors))
		for k2, v2 := range o.RegistryMirrors {
			var cp_RegistryMirrors_v2 *RegistryMirrorConfig
			if v2 != nil {
				cp_RegistryMirrors_v2 = new(RegistryMirrorConfig)
				*cp_RegistryMirrors_v2 = *v2
				if v2.MirrorEndpoints != nil {
					cp_RegistryMirrors_v2.MirrorEndpoints = make([]RegistryEndpointConfig, len(v2.MirrorEndpoints))
					copy(cp_RegistryMirrors_v2.MirrorEndpoints, v2.MirrorEndpoints)
				}
				if v2.MirrorSkipFallback != nil {
					cp_RegistryMirrors_v2.MirrorSkipFallback = new(bool)
					*cp_RegistryMirrors_v2.MirrorSkipFallback = *v2.MirrorSkipFallback
				}
			}
			cp.RegistryMirrors[k2] = cp_RegistryMirrors_v2
		}
	}
	if o.RegistryConfig != nil {
		cp.RegistryConfig = make(map[string]*RegistryConfig, len(o.RegistryConfig))
		for k2, v2 := range o.RegistryConfig {
			var cp_RegistryConfig_v2 *RegistryConfig
			if v2 != nil {
				cp_RegistryConfig_v2 = new(RegistryConfig)
				*cp_RegistryConfig_v2 = *v2
				if v2.RegistryTLS != nil {
					cp_RegistryConfig_v2.RegistryTLS = new(RegistryTLSConfig)
					*cp_RegistryConfig_v2.RegistryTLS = *v2.RegistryTLS
					if v2.RegistryTLS.TLSClientIdentity != nil {
						cp_RegistryConfig_v2.RegistryTLS.TLSClientIdentity = v2.RegistryTLS.TLSClientIdentity.DeepCopy()
					}
					cp_RegistryConfig_v2.RegistryTLS.TLSCA = v2.RegistryTLS.TLSCA.DeepCopy()
					if v2.RegistryTLS.TLSInsecureSkipVerify != nil {
						cp_RegistryConfig_v2.RegistryTLS.TLSInsecureSkipVerify = new(bool)
						*cp_RegistryConfig_v2.RegistryTLS.TLSInsecureSkipVerify = *v2.RegistryTLS.TLSInsecureSkipVerify
					}
				}
				if v2.RegistryAuth != nil {
					cp_RegistryConfig_v2.RegistryAuth = new(RegistryAuthConfig)
					*cp_RegistryConfig_v2.RegistryAuth = *v2.RegistryAuth
				}
			}
			cp.RegistryConfig[k2] = cp_RegistryConfig_v2
		}
	}
	return cp
}

// DeepCopy generates a deep copy of ImageCacheConfigSpec.
func (o ImageCacheConfigSpec) DeepCopy() ImageCacheConfigSpec {
	var cp ImageCacheConfigSpec = o
	if o.Roots != nil {
		cp.Roots = make([]string, len(o.Roots))
		copy(cp.Roots, o.Roots)
	}
	return cp
}

// DeepCopy generates a deep copy of SeccompProfileSpec.
func (o SeccompProfileSpec) DeepCopy() SeccompProfileSpec {
	var cp SeccompProfileSpec = o
	if o.Value != nil {
		cp.Value = make(map[string]any, len(o.Value))
		for k2, v2 := range o.Value {
			cp.Value[k2] = v2
		}
	}
	return cp
}
