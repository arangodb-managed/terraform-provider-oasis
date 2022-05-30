//
// DISCLAIMER
//
// Copyright 2020-2022 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//

package internal

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

// counter is keeping track of the generated resources in an atomic way.
// This will result in unique ids even with multiple terraform calls.
var counter uint64

func uniqueResourceID(prefix string) string {
	ts := strings.ReplaceAll(time.Now().Format("200601020405.000000"), ".", "")
	atomic.AddUint64(&counter, 1)
	return fmt.Sprintf("%s%s%05d", prefix, ts, counter)
}
