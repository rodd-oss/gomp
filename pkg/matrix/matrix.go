/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

none :)

Thank you for your support!
*/

package matrix

type V4 [4]float32
type M4 [16]float32

//go:nosplit
func multiplyasm(data []V4, m M4)

func M4MultiplyV4(m M4, v V4) V4 {
	return V4{
		v[0]*m[0] + v[1]*m[4] + v[2]*m[8] + v[3]*m[12],
		v[0]*m[1] + v[1]*m[5] + v[2]*m[9] + v[3]*m[13],
		v[0]*m[2] + v[1]*m[6] + v[2]*m[10] + v[3]*m[14],
		v[0]*m[3] + v[1]*m[7] + v[2]*m[11] + v[3]*m[15],
	}
}

func multiply(data []V4, m M4) {
	for i, v := range data {
		data[i] = M4MultiplyV4(m, v)
	}
}
