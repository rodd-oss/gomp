// func multiplyasm(data []V4, m M4)
//
// memory layout of the stack relative to FP
//  +0 data slice ptr
//  +8 data slice len
// +16 data slice cap
// +24 m[0]  | m[1]
// +32 m[2]  | m[3]
// +40 m[4]  | m[5]
// +48 m[6]  | m[7]
// +56 m[8]  | m[9]
// +64 m[10] | m[11]
// +72 m[12] | m[13]
// +80 m[14] | m[15]

#include "textflag.h"

TEXT ·multiplyasm(SB),NOSPLIT,$0
  // data ptr
  MOVQ data+0(FP), CX
  // data len
  MOVQ data+8(FP), SI
  // index into data
  MOVQ $0, AX
  // return early if zero length
  CMPQ AX, SI
  JE END
  // load the matrix into 128-bit wide xmm registers
  // load [m[0], m[1], m[2], m[3]] into xmm0
  MOVUPS m+24(FP), X0
  // load [m[4], m[5], m[6], m[7]] into xmm1
  MOVUPS m+40(FP), X1
  // load [m[8], m[9], m[10], m[11]] into xmm2
  MOVUPS m+56(FP), X2
  // load [m[12], m[13], m[14], m[15]] into xmm3
  MOVUPS m+72(FP), X3
LOOP:
  // load each component of the vector into xmm registers
  // load data[i][0] (x) into xmm4
  MOVSS    0(CX), X4
  // load data[i][1] (y) into xmm5
  MOVSS    4(CX), X5
  // load data[i][2] (z) into xmm6
  MOVSS    8(CX), X6
  // load data[i][3] (w) into xmm7
  MOVSS    12(CX), X7
  // copy each component of the matrix across each register
  // [0, 0, 0, x] => [x, x, x, x]
  SHUFPS $0, X4, X4
  // [0, 0, 0, y] => [y, y, y, y]
  SHUFPS $0, X5, X5
  // [0, 0, 0, z] => [z, z, z, z]
  SHUFPS $0, X6, X6
  // [0, 0, 0, w] => [w, w, w, w]
  SHUFPS $0, X7, X7
  // xmm4 = [m[0], m[1], m[2], m[3]] .* [x, x, x, x]
  MULPS X0, X4
  // xmm5 = [m[4], m[5], m[6], m[7]] .* [y, y, y, y]
  MULPS X1, X5
  // xmm6 = [m[8], m[9], m[10], m[11]] .* [z, z, z, z]
  MULPS X2, X6
  // xmm7 = [m[12], m[13], m[14], m[15]] .* [w, w, w, w]
  MULPS X3, X7
  // xmm4 = xmm4 + xmm5
  ADDPS X5, X4
  // xmm4 = xmm4 + xmm6
  ADDPS X6, X4
  // xmm4 = xmm4 + xmm7
  ADDPS X7, X4
  // data[i] = xmm4
  MOVNTPS X4, 0(CX)
  // data++
  ADDQ $16, CX
  // i++
  INCQ AX
  // if i >= len(data) break
  CMPQ AX, SI
  JLT LOOP
END:
  // since we use a non-temporal write (MOVNTPS)
  // make sure all writes are visible before we leave the function
  SFENCE
  RET
