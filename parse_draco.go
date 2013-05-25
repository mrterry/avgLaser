package goAvgLaser

import "io/ioutil"
import "bytes"
import "strconv"


func Parse(path string, r0 float64) (pointings []Pointing, rings [][]int, spots []PDDSpot, interp int) {
  text, err := ioutil.ReadFile(path)
  if err != nil {
    println("Broken")
    panic(err)
  }
  text = strip_comments(text)
  words := bytes.Fields(text)

  n_segs, words := take_ints(1, words)
  ms, words := take_ints(n_segs[0], words)
  pointings, words = take_pointings(ms[0], words)

  end := 5 + ms[n_segs[0]-1] - 1
  rings, words = take_rings(ms[4:end], words)
  interp = 0
  spots, words = take_spots(ms[len(ms)-1], words)

  if len(spots) != len(rings) {
    panic(-1)
  }
  for i := 0; i < len(spots); i++ {
    spots[i].Radius *= r0
    spots[i].SetOff1(pointings[rings[i][0]])
  }
  return
}


func strip_comments(s []byte) ([]byte) {
  for {
    i := bytes.Index(s, []byte("!"))
    if i == -1 {
      break
    }
    j := bytes.Index(s, []byte("\n"))
    if j == -1 {
      println("bad!")
      panic(j)
    }
    s = s[j+1:]
  }
  return s
}


func take_ints(n int, words [][]byte) (ints []int, words_left[][]byte) {
  ints = make([]int, n)
  for i, w := range words[:n] {
    j, err := strconv.Atoi(string(w))
    if err != nil {
      panic(err)
    }
    ints[i] = j
  }
  words_left = words[n:]
  return
}


func take_float64s(n int, words [][]byte) (floats []float64, words_left [][]byte) {
  floats = make([]float64, n)
  for i, w := range words[:n] {
    f, err := strconv.ParseFloat(string(w), 64)
    if err != nil {
      panic(err)
    }
    floats[i] = f
  }
  words_left = words[n:]
  return
}


func take_pointings(n int, words [][]byte) (pointings []Pointing, words_left [][]byte) {
  pointings = make([]Pointing, n)

  thetas, words := take_float64s(n, words)
  phis, words := take_float64s(n, words)
  theta_rps, words := take_float64s(n, words)
  phi_rps, words := take_float64s(n, words)

  for i := 0; i < n; i++ {
    pointings[i].Theta = thetas[i]
    pointings[i].Phi = phis[i]
    pointings[i].Theta_rp = theta_rps[i]
    pointings[i].Phi_rp = phi_rps[i]
  }
  words_left = words
  return
}


func take_rings(ring_beams []int, words [][]byte) (rings [][]int, words_left [][]byte) {
  words_left = words
  n_rings := len(ring_beams)
  rings = make([][]int, n_rings)
  for i, n_beams := range ring_beams {
    rings[i], words_left = take_ints(n_beams, words_left)
  }
  return
}


func take_spots(n_spots int, words [][]byte) (spots []PDDSpot, words_left [][]byte) {
  spots = make([]PDDSpot, n_spots)

  words_left = words
  sg1, words_left := take_float64s(n_spots, words_left)
  radius, words_left := take_float64s(n_spots, words_left)
  ellip1, words_left := take_float64s(n_spots, words_left)

  amp2, words_left := take_float64s(n_spots, words_left)
  sg2, words_left := take_float64s(n_spots, words_left)
  off2, words_left := take_float64s(n_spots, words_left)
  ellip2, words_left := take_float64s(n_spots, words_left)
  
  for i := 0; i < len(spots); i++ {
    spots[i].Radius = radius[i]
    spots[i].SG1 = sg1[i]
    spots[i].SG2 = sg2[i]
    spots[i].Ellip1 = ellip1[i]
    spots[i].Ellip2 = ellip2[i]
    spots[i].Amp1 = 1.
    spots[i].Amp2 = amp2[i]
    // Off1
    spots[i].Off2 = off2[i]
  }
  return
}
