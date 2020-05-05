// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

<<<<<<< HEAD
// At this point only s390x has an assembly implementation of sum. All other
// platforms have assembly implementations of mac, and just define sum as using
// that through New. Once s390x is ported, this file can be deleted and the body
// of sum moved into Sum.

// +build !go1.11 !s390x gccgo purego
=======
// +build s390x,!go1.11 !amd64,!s390x,!ppc64le gccgo purego
>>>>>>> origin/r_permissions

package poly1305

func sum(out *[TagSize]byte, msg []byte, key *[32]byte) {
<<<<<<< HEAD
	h := New(key)
=======
	h := newMAC(key)
>>>>>>> origin/r_permissions
	h.Write(msg)
	h.Sum(out[:0])
}
