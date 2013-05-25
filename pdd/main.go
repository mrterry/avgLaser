package main

import al "github.com/mrterry/goAvgLaser"
import "math"

func main() {
  r0 := 1485.e-4
  threads := 1
  nbins := 100
  lens_hw := 40./2
  chunk_size := 1000
  nrays := chunk_size*100

  tally_edges := al.Linspace(0., math.Pi, nbins+1)

  wall_pattern := al.WallPattern{lens_hw, lens_hw, 500}
  focus_edges := al.Linspace(-0.2, 0.2, 200)

  transpChan := make(chan al.TransportJob, threads*2)
  histChan := make(chan []int, threads)
  for i:=0; i<threads; i++ {
    go al.TransportAndTally(&wall_pattern, tally_edges, chunk_size, r0, transpChan, histChan)
  }

  points, rings, spots, _ := al.Parse("/Users/terry10/code/average_laser/NIFPortConfig.dat", r0)
  for spot_index, spot := range spots {
    focus_pattern := al.NewSpotPattern(&spot, focus_edges, focus_edges)
    for _, beam_index := range rings[spot_index] {
      transpChan <- al.TransportJob{&focus_pattern, points[beam_index-1], nrays}
    }
  }
  close(transpChan)

  total_tally := make([]int, nbins)
  for i:=0; i<threads; i++ {
    local_tally := <- histChan
    for n, count := range local_tally {
      total_tally[n] += count
    }
  }
  close(histChan)

  println(total_tally)
}
