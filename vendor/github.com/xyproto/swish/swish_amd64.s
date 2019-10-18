#include "textflag.h"

DATA expodata<>+0(SB)/8, $1.0
DATA expodata<>+8(SB)/8, $-0.00390625
GLOBL expodata<>+0(SB), RODATA, $16

TEXT ·SwishAssembly(SB),NOSPLIT|NOPTR,$0-16
  // x+0(FP) is the given argument
  MOVSD x+0(FP), X1
  MOVSD x+0(FP), X2
  // x1 *= -0.00390625 which is (1/256)
  MULSD expodata<>+8(SB), X1
  // x1 += 1.0
  MOVSD expodata<>+0(SB), X3
  ADDSD X3, X1
  // x1 *= x1 ...
  MULSD X1, X1
  MULSD X1, X1
  MULSD X1, X1
  MULSD X1, X1
  MULSD X1, X1
  MULSD X1, X1
  MULSD X1, X1
  MULSD X1, X1
  // x1 += 1.0
  ADDSD X3, X1
  // x2 /= x1
  DIVSD X1, X2
  // return x2
  MOVSD X2, ret+8(FP)
  // done, jump back
  RET
