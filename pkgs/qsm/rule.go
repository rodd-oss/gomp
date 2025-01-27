/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package qsm

import "time"

type MutationBeforeHandler = func()
type MutationUpdateHandler = func(dt time.Duration)
type MutationAfterHandler = func()
type MutationCancelHandler = func()

type MutationRule struct {
	Before MutationBeforeHandler
	While  MutationUpdateHandler
	After  MutationAfterHandler
	Cancel MutationCancelHandler

	Duration time.Duration
}

func (t *MutationRule) runBefore() {
	f := t.Before
	if f != nil {
		f()
	}
}

func (t *MutationRule) runWhile(dt time.Duration) {
	f := t.While
	if f != nil {
		f(dt)
	}
}

func (t *MutationRule) runAfter() {
	f := t.After
	if f != nil {
		f()
	}
}

func (t *MutationRule) runCancel() {
	if t == nil {
		return
	}

	f := t.Cancel
	if f != nil {
		f()
	}
}
