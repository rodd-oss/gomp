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

package worker

type AnyTask interface {
	Run(workerId WorkerId) error
}

type TaskError struct {
	Err error
	Id  WorkerId
}

func (e TaskError) Error() string {
	return e.Err.Error()
}
