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

package timeasm

import (
	"math/bits"
	"time"
	_ "unsafe" // Для go:linkname
)

var (
	tscFrequency      uint64
	baseTimeUnix      int64
	tscFrequencyShift int
)

//func init() {
//	calibrateTSC()
//}

func calibrateTSC() {
	baseTime := time.Now()
	baseTimeUnix = baseTime.UnixNano()
	ticksStart := Cputicks()
	time.Sleep(100 * time.Millisecond)
	endTime := time.Now()
	ticksEnd := Cputicks()

	tscFrequency = (ticksEnd - ticksStart) * 10 / uint64(endTime.Sub(baseTime).Nanoseconds())
	tscFrequencyShift = bits.Len64(tscFrequency) - 1
}

func GetTimestamp() uint64 {
	return Cputicks()
}

var baseNano int64 = 0
var baseTsc uint64
var scale float64

func Now() int64 {
	if baseNano == 0 {
		baseNano = time.Now().UnixNano()
		baseTsc = Cputicks()
		time.Sleep(1e7)
		rn := time.Now().UnixNano()
		ct := Cputicks()
		scale = float64(rn-baseNano) / float64(ct-baseTsc)
	}
	return baseNano + int64(float64(Cputicks()-baseTsc)*scale)
}
