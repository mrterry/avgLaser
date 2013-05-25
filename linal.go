package avgLaser

import "math"
import "errors"
import "sort"


func Linspace(begin, end float64, num int) (x []float64) {
  x = make([]float64, num)
  dx := (end - begin)/float64(num-1)
  for i:=0; i<int(num)-1; i++ {
    x[i] = begin + float64(i)*dx
  }
  x[num-1] = end
  return
}


func Meshgrid(x, y []float64) (X, Y Cont2DArray) {
  ni, nj := len(x), len(y)
  X = NewCont2DArray(ni, nj)
  Y = NewCont2DArray(ni, nj)
  j := 0
  for i:=0; i<ni; i++ {
    X[i][j] = x[i]
    for j=0; j<nj; j++ {
      Y[i][j] = y[j]
    }
  }
  return
}


type Cont2DArray [][]float64
func (a Cont2DArray) Flat() (b []float64) {
  ni, nj := len(a), len(a[0])
  b = a[0][:ni*nj]
  return
}
func (a *Cont2DArray) Reshape2(ni, nj int) (c Cont2DArray) {
  flat := a.Flat()
  return Reshape2(flat, ni, nj)
}
func NewCont2DArray(ni, nj int) (x Cont2DArray){
  data := make([]float64, ni*nj)
  x = Reshape2(data, ni, nj)
  return
}


func Reshape2(data []float64, ni, nj int) (y Cont2DArray) {
  x := make([][]float64, ni)
  for i:=0; i<ni; i++ {
    x[i], data = data[:nj], data[nj:]
  }
  y = Cont2DArray(x)
  return
}


func StableQuad(a, b, c float64) (l1, l2 float64, err error) {
  q := b*b - 4*a*c
  if q < 0. {
    err = errors.New("complex roots")
    return
  }
  q = -0.5 * (b + math.Copysign(math.Sqrt(q), b))
  l1, l2 = q/a, c/q
  return
}


func Histogram0(data, edges []float64, tally []int) (n_lo, n_hi int) {
  // Sort then histogram
  // Use if you have a smallish number of points and large-ish number
  // of bins
  sort.Float64s(data)
  n_edges := len(edges)

  var x float64
  n_lo = 0
  bot := edges[0]
  for _, x = range data {
    if x < bot {
      n_lo += 1
    } else {
      break
    }
  }

  e := 1
  for i, x := range data[n_lo:] {
    for ; e < n_edges; e++ {
      if x <= edges[e] {
        tally[e-1] += 1
        break
      }
    } 
    if e == n_edges {
      // if reach the last bin, then everything else is hi
      n_hi = len(data) - n_lo - i
      return
    }
  }
  return
}


func Histogram1(data, edges []float64, tally []int) (n_lo, n_hi int) {
  // Brute force search
  // Use this one, generally the fastest, unles you have a huge number of bins
  ne := len(edges)
  bot := edges[0]
  top := edges[ne-1]

  for _, x := range data {
    if x < bot {
      n_lo += 1
      continue
    } 
    if x > top {
      n_hi += 1
      continue
    }

    for b, e := range edges[1:] {
      if x < e {
        tally[b] += 1
        break
      }
    }
  }
  return
}


func Histogram2(data, edges []float64, tally []int) (n_lo, n_hi int) {
  // Binary search
  // Use this one if you have a very large number of bins
  ne := len(edges)

  for _, x := range data {
    i := sort.SearchFloat64s(edges, x)
    if i == 0 {
      n_lo += 1
    } else if i == ne {
      n_hi += 1
    } else {
      tally[i-1] += 1
    }
  }
  return
}


func Zcen(x [][]float64) (y Cont2DArray) {
  ni, nj := len(x), len(x[0])

  y = NewCont2DArray(ni-1, nj-1)

  for i, row := range y {
    for j := range row {
      y[i][j] = 0.25 * (x[i][j] + x[i+1][j] + x[i][j+1] + x[i+1][j+1])
    }
  }
  return
}


func CumSum(x []float64) (y []float64) {
  ny := len(x)
  y = make([]float64, ny)

  for i, val := range x {
    for j:=i; j<ny; j++ {
      y[j] += val
    }
  }
  return
}
