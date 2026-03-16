/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package ecs

import "math/bits"

func FastIntLog2(value uint64) int {
	return bits.Len64(value) - 1
}

func FastestIntLog2(value uint64) int {
	return bits.LeadingZeros64(value) ^ 63
}
