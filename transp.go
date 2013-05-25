package avgLaser

import "code.google.com/p/biogo.blas"
import "math"
import "errors"


func GetMu(p, q []float64, r float64) (mu float64, err error) {
  s, R2 := 0., 0.
  for i := range q {
    s += p[i]*q[i]
    R2 += p[i]*p[i]
  }

  c := R2 - r*r
  d := s*s - c
  if d < 0. {
    err = errors.New("complex roots")
    return
  }
  d = math.Sqrt(d) - s

  l1, l2 := d, c/d

  if err != nil {
    return
  }
  mu = -(s + math.Min(l1, l2))/r
  return
}


func RotateThetaPhi(points Cont2DArray, theta, phi float64) {
  n, d := len(points), len(points[0])

  ct, st := math.Cos(theta), math.Sin(theta)
  cp, sp := math.Cos(phi), math.Sin(phi)
  x := points.Flat()
  y, z := x[1:], x[2:]

  param := new(blas.DrotmParams)
  param.Flag = -1.
  param.H[0] = -st
  param.H[1] = ct
  param.H[2] = st
  param.H[3] = ct
  blas.Drotm(n, x, d, z, d, param)

  param.H[0] = cp
  param.H[1] = sp
  param.H[2] = -cp
  param.H[3] = sp
  blas.Drotm(n, x, d, y, d, param)
}


func Focus2Q(walls, foci Cont2DArray) {
  tmp := 0.
  for i, q := range foci {
    norm := 0.
    for j, v := range q {
      tmp = v - walls[i][j]
      foci[i][j] = tmp
      norm += tmp*tmp
    }
    norm = 1/math.Sqrt(norm)
    for j := range q {
      foci[i][j] *= norm
    }
  }
  return
}
