
using import "EnumEnumCityEnumParks.capnp".EnumEnumCityEnumParks;
@0x851c6a9c2bfa861e;

struct EnumCity {
  enumParks @0 :EnumEnumCityEnumParks;
  name @1 :Text;
}