/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import "math/bits"

func FastIntLog2(value int) int {
	return bits.Len64(uint64(value)) - 1
}
