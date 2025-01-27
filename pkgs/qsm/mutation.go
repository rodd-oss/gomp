/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package qsm

type Mutation[T State] struct {
	targetState *T
	progress    float32 // 0.0 - 1.0
	canceler    chan bool
	rule        *MutationRule
}

func (t *Mutation[T]) reset() {
	t.targetState = nil
	t.progress = 0.0
	t.canceler = make(chan bool, 1)
	t.rule = nil
}
