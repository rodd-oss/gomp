/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package components

import (
	"math"

	"capnproto.org/go/capnp/v3"
	"github.com/jakecoffman/cp/v2"
	ecs "github.com/yohamta/donburi"
	ecsmath "github.com/yohamta/donburi/features/math"
)

type PhysicsData struct {
	Body *cp.Body
}

const (
	interpolationSpeed = 0.25
)

func (p PhysicsData) Update(dt float64, e *ecs.Entry, isClient bool) error {
	if p.Body.IsSleeping() {
		return nil
	}

	pos := p.Body.Position()
	pos.X = math.Round(pos.X)
	pos.Y = math.Round(pos.Y)
	// round body position to nearest integer
	p.Body.SetPosition(pos.Lerp(pos, interpolationSpeed))

	lastTransformPosition := Transform.GetValue(e).LocalPosition
	newTransformPosition := ecsmath.NewVec2(pos.X, pos.Y)

	posDelta := &ecsmath.Vec2{
		X: p.Body.Position().X - lastTransformPosition.X,
		Y: p.Body.Position().Y - lastTransformPosition.Y,
	}

	if isClient {
		Transform.SetValue(e, TransformData{
			LocalPosition: newTransformPosition.Add(posDelta.MulScalar(-0.66)),
		})
	} else {
		Transform.SetValue(e, TransformData{
			LocalPosition: newTransformPosition,
		})
	}

	return nil
}

func (p PhysicsData) OnSerialize(e *ecs.Entry) {
	arena := capnp.SingleSegment(nil)

	// Make a brand new empty message.  A Message allocates Cap'n Proto structs within
	// its arena.  For convenience, NewMessage also returns the root "segment" of the
	// message, which is needed to instantiate the Book struct.  You don't need to
	// understand segments and roots yet (or maybe ever), but if you're curious, messages
	// and segments are documented here:  https://capnproto.org/encoding.html
	msg, seg, err := capnp.NewMessage(arena)
	if err != nil {
		panic(err)
	}

	// Create a new Book struct.  Every message must have a root struct.  Again, it is
	// not important to understand "root structs" at this point.  For now, just understand
	// that every type you instantiate needs to be a "root", unless you plan on assigning
	// it to another object.  When in doubt, use NewRootXXX.
	//
	// If you're insatiably curious, see:  https://capnproto.org/encoding.html#messages
	book, err := books.NewRootBook(seg)
	if err != nil {
		panic(err)
	}

	// Great, we have our book!  Now let's set some fields.  Each field you declared in
	// your schema will produce two methods on the generated type.  The "getter" method
	// has the name of the field, for example:  Book.Title().  The corresponding "setter"
	// method is prefixed with "Set", for example:  Book.SetTitle().
	//
	// Some getters and setters return errors, which we are ignoring in this example for
	// the sake of clarity.  Your code SHOULD check these errors and handle them.
	//
	// To begin, we set the book's title to "War and Peace".
	_ = book.SetTitle("War and Peace")

	// Then, we set the page count.
	book.SetPageCount(1440)
}

var Physics = ecs.NewComponentType[PhysicsData]()
