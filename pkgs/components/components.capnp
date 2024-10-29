using Go = import "/go.capnp";
using Physics = import "./physics/physics.capnp";
@0xabb220532a204a11;
$Go.package("components");
$Go.import("components");

struct ComponentState {
  state :union {
    none @0 :Void;
    physics @1 :Physics.PhysicsState;
  }
}

struct ComponentsStates {
  states @0 :List(ComponentState);
}