using Go = import "/go.capnp";
@0xa61827898c3c5ebe;
$Go.package("physics");
$Go.import("gomp_game/pkgs/components/physics");

struct PhysicsState {
  x @0 :Int32;
  y @1 :Int32;
}

