// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
/**
Provide a mechanism to compute the checksum of a rolling window
in single byte increments by using the hash parts computed earlier
**/

package xferspdy

import (
	"github.com/golang/glog"
)

const (
	// mod is the largest prime that is less than 65536.
	mod = 65521
	//number of bytes that can be added
	nmax = 5552
)

// The low 16 bits are s1, the high 16 bits are s2.
type checksum uint32

//State of Adler-32 computation
//It contants, the byte arary window from the most recent computation
//and interim sum values
type State struct {
	window []byte
	s1     uint32
	s2     uint32
}

//Checksum returns the Adler-32 checksum, computed for the given byte array
//In addition, returns a State that captures the interim results during computation
//The State can be used to update the byte[] window and compute rolling hash
func Checksum(p []byte) (uint32, *State) {
	glog.V(4).Infof("Length of buffer %d \n Calc checksum for \n %v \n", len(p), p)
	s1, s2 := uint32(1&0xffff), uint32(1>>16)
	glog.V(4).Infof("Init: s1 %d s2 %d\n", s1, s2)
	orig := p
	for len(p) > 0 {
		var q []byte
		if len(p) > nmax {
			p, q = p[:nmax], p[nmax:]
		}
		for _, x := range p {
			s1 += uint32(x)
			s2 += s1
		}
		s1 %= mod
		s2 %= mod
		p = q
	}
	glog.V(4).Infof("s1 %d s2 %d\n", s1, s2)
	return uint32(s2<<16 | s1), &State{orig, s1, s2}
}

//UpdateWindow returns rolling Adler-32 checksum after updating the window
//The checksum is not calculated from scratch.
//Instead the captured byte array window in State struct is updated,
//similar to circular buffer and rolling hash is calculated
func (s *State) UpdateWindow(nb byte) uint32 {
	s.window = append(s.window, nb)
	x := s.window[0]
	s.window = s.window[1:]
	s.s1 = s.s1 + uint32(nb) - uint32(x)
	s.s1 %= mod
	b := (uint32(len(s.window)) * uint32(x)) + 1
	a := s.s2 + s.s1
	for b > a {
		a += mod
	}
	s.s2 = a - b
	s.s2 %= mod
	return uint32(s.s2<<16 | s.s1)
}