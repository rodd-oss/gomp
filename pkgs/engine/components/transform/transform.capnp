@0xe345bdf5f8f7863c;
using Go = import "/go.capnp";

$Go.package("transform");
$Go.import("gomp_game/pkgs/components/transform");

struct TransformState {
  x @0 :Int32;
  y @1 :Int32;
}