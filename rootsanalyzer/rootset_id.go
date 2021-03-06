// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rootsanalyzer

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/google/certificate-transparency-go/x509"
	"github.com/google/monologue/storage"
)

// GenerateCertID returns SHA-256 of certificate's DER representation.
func GenerateCertID(root *x509.Certificate) ([32]byte, error) {
	if root == nil {
		return [32]byte{}, fmt.Errorf("unable to generate root-ID for nil")
	}
	return sha256.Sum256(root.Raw), nil
}

// GenerateSetID returns ID for a set of root-certificates.
// Order or multi-entries of a same certificate do not influence the ID.
func GenerateSetID(roots []*x509.Certificate) (storage.RootSetID, error) {
	var dedupRootIDs []string
	rootIDSet := make(map[[32]byte]bool)
	for _, r := range roots {
		certID, err := GenerateCertID(r)
		if err != nil {
			return "", fmt.Errorf("unable to generate ID for certificate %q: %s", r.Subject, err)
		}
		if !(rootIDSet[certID]) {
			rootIDSet[certID] = true
			dedupRootIDs = append(dedupRootIDs, string(certID[:]))
		}
	}

	// sort deduplicated roots IDs
	sort.Strings(dedupRootIDs)

	// concatenate roots IDs
	concat := strings.Join(dedupRootIDs[:], "")
	concatHash := sha256.Sum256([]byte(concat))
	return storage.RootSetID(string(concatHash[:])), nil
}
