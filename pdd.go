package goAvgLaser

import "sort"
import "math/rand"


type Sampler interface {
  SampleTo(Cont2DArray)
}


type Pointing struct {
  Theta, Phi, Theta_rp, Phi_rp float64
}


type Spotish interface {
  Intensity(x, y []float64) [][]float64
  GetRadius() float64
}


type TransportJob struct {
  Pattern Sampler
  P Pointing
  NRays int
}


type WallPattern struct {
  LensWidth, LensHeight, WallRadius float64
}
func (wp *WallPattern) SampleTo(wall_points Cont2DArray) {
  for i := range wall_points {
    wall_points[i][0] = wp.LensWidth*(rand.Float64() - 0.5)
    wall_points[i][0] = wp.LensHeight*(rand.Float64() - 0.5)
    wall_points[i][2] = wp.WallRadius
  }
  return
}


type SpotPattern struct {
  spot Spotish
  edges_x, edges_y []float64
  cum_prob []float64
}
func (sp *SpotPattern) SampleTo(spot_points Cont2DArray) {
  // If not already created, create cumulative probability distribution
  if sp.cum_prob == nil {
    intensity := sp.spot.Intensity(sp.edges_x, sp.edges_y)
    cell_power := Zcen(intensity)
    for i, row := range cell_power {
      dx := sp.edges_x[i+1] - sp.edges_x[i]
      for j, _ := range row {
        dy := sp.edges_y[i+1] - sp.edges_y[i]
        cell_power[i][j] *= dx*dy
      }
    }

    // Normalize to 1
    sp.cum_prob = CumSum(cell_power.Flat())
    max_inv := 1./sp.cum_prob[len(cell_power)-1]
    for i := range(sp.cum_prob) {
      sp.cum_prob[i] *= max_inv
    }
  }

  ny := len(sp.edges_y) - 1
  for s := range spot_points {
    samp := rand.Float64()
    k := sort.SearchFloat64s(sp.cum_prob, samp)
    i, j := k/ny, k%ny
    dx := sp.edges_x[i+1] - sp.edges_x[i]
    dy := sp.edges_y[j+1] - sp.edges_y[j]
    spot_points[s][0] = sp.edges_x[i] + rand.Float64()*dx
    spot_points[s][1] = sp.edges_y[j] + rand.Float64()*dy
    spot_points[s][2] = sp.spot.GetRadius()
  }
  return
}
func NewSpotPattern(spot Spotish, edges_x, edges_y []float64) (sp SpotPattern) {
  sp = SpotPattern{spot, edges_x, edges_y, nil}
  return
}


func Tally(edges []float64, muChan chan []float64, doneChan chan []int) {
  hist := make([]int, len(edges)-1)
  for mus := range muChan {
    Histogram1(mus, edges, hist)
  }
  doneChan <- hist
  return
}


func TransportAndTally(wall_pattern Sampler, edges []float64, chunk int,
    target_radius float64, transpChan chan TransportJob, histChan chan []int) {
  walls := NewCont2DArray(chunk, 3)
  foci := NewCont2DArray(chunk, 3)
  mus := make([]float64, 0, chunk)
  hist := make([]int, len(edges)-1)

  npieces := 0
  for job := range transpChan {
    left := job.NRays
    for left > chunk {
      if left > chunk {
        walls = walls[:chunk]
        foci = foci[:chunk]
        left -= chunk
      } else {
        walls = walls[:left]
        foci = foci[:left]
        left = 0
      }
      npieces += 1

      wall_pattern.SampleTo(walls)
      RotateThetaPhi(walls, job.P.Theta, job.P.Phi)

      job.Pattern.SampleTo(foci)
      RotateThetaPhi(foci, job.P.Theta_rp, job.P.Phi_rp)

      Focus2Q(walls, foci)
      qs := foci

      // TODO: Race condition here
      // tabChan may not be done with mus before I reuse it
      mus = mus[:0]
      for i := range qs {
        mu, err := GetMu(walls[i], qs[i], target_radius)
        if err != nil {
          continue
        }
        mus = append(mus, mu)
      }
      Histogram1(mus, edges, hist)
    }
  }
  println(&hist)
  histChan <- hist
}
