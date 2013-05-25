package goAvgLaser

import "math"


type PDDSpot struct {
  Radius, SG1, SG2, Ellip1, Ellip2, Amp1, Amp2, Off1, Off2 float64
}
func (s *PDDSpot) SetOff1(p Pointing) {
  s.Off1 = (0.5 * s.Radius *
      math.Sin(math.Pi/180 * (p.Theta_rp-p.Theta)) / math.Sin(p.Theta))
  return
}
func (s *PDDSpot) Intensity(x, y []float64) (I [][]float64){
  ni := len(x)
  nj := len(y)
  a := math.Log(0.05)

  sg_apr := 10.
  r_apr := 1.15*s.Radius

  I = NewCont2DArray(ni, ni)
  I_max := 0.
  for i:=0; i<ni; i++ {
    x2 := math.Pow(x[i], 2)
    for j:=0; j<nj; j++ {
      my_y := y[j] 
      R0 := x2 + math.Pow(my_y - s.Off1, 2)
      R1 := x2 + math.Pow(my_y, 2)
      R2 := x2 + math.Pow(s.Ellip2*(my_y - s.Off2), 2)

      I[i][j] += s.Amp1 * math.Exp(-a*math.Pow(R1/s.Radius, s.SG1))
      I[i][j] += s.Amp2 * math.Exp(-a*math.Pow(R2/s.Radius, s.SG2))
      I[i][j] *= math.Exp(-a * math.Pow(R0/r_apr, sg_apr))
      I_max = math.Max(I_max, I[i][j])
    }
  }

  for i:=0; i<ni; i++ {
    for j:=0; j<nj; j++ {
      I[i][j] /= I_max
    }
  }
  return
}
func (s *PDDSpot) GetRadius() (float64) {
  return s.Radius
}
