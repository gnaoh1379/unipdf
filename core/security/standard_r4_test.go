/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package security

import (
	"github.com/gnaoh1379/unipdf/v3/common"
	"testing"
)

func init() {
	common.SetLogger(common.ConsoleLogger{})
}

func TestR4Padding(t *testing.T) {
	sh := stdHandlerR4{}

	// Case 1 empty pass, should match padded string.
	key := sh.paddedPass([]byte(""))
	if len(key) != 32 {
		t.Errorf("Fail, expected padded pass length = 32 (%d)", len(key))
	}
	if key[0] != 0x28 {
		t.Errorf("key[0] != 0x28 (%q in %q)", key[0], key)
	}
	if key[31] != 0x7A {
		t.Errorf("key[31] != 0x7A (%q in %q)", key[31], key)
	}

	// Case 2, non empty pass.
	key = sh.paddedPass([]byte("bla"))
	if len(key) != 32 {
		t.Errorf("Fail, expected padded pass length = 32 (%d)", len(key))
	}
	if string(key[0:3]) != "bla" {
		t.Errorf("Expecting start with bla (%s)", key)
	}
	if key[3] != 0x28 {
		t.Errorf("key[3] != 0x28 (%q in %q)", key[3], key)
	}
	if key[31] != 0x64 {
		t.Errorf("key[31] != 0x64 (%q in %q)", key[31], key)
	}
}

// Test algorithm 2.
func TestAlg2(t *testing.T) {
	sh := stdHandlerR4{
		// V: 2,
		ID0: string([]byte{0x4e, 0x00, 0x99, 0xe5, 0x36, 0x78, 0x93, 0x24,
			0xff, 0xd5, 0x82, 0xe4, 0xec, 0x0e, 0xa3, 0xb4}),
		Length: 128,
	}
	d := &StdEncryptDict{
		R:               3,
		P:               0xfffff0c0,
		EncryptMetadata: true,
		O: []byte{0xE6, 0x00, 0xEC, 0xC2, 0x02, 0x88, 0xAD, 0x8B,
			0x5C, 0x72, 0x64, 0xA9, 0x5C, 0x29, 0xC6, 0xA8, 0x3E, 0xE2, 0x51,
			0x76, 0x79, 0xAA, 0x02, 0x18, 0xBE, 0xCE, 0xEA, 0x8B, 0x79, 0x86,
			0x72, 0x6A, 0x8C, 0xDB},
	}

	key := sh.alg2(d, []byte(""))

	keyExp := []byte{0xf8, 0x94, 0x9c, 0x5a, 0xf5, 0xa0, 0xc0, 0xca,
		0x30, 0xb8, 0x91, 0xc1, 0xbb, 0x2c, 0x4f, 0xf5}

	if string(key) != string(keyExp) {
		common.Log.Debug("   Key (%d): % x", len(key), key)
		common.Log.Debug("KeyExp (%d): % x", len(keyExp), keyExp)
		t.Errorf("alg2 -> key != expected\n")
	}
}

// Test algorithm 3.
func TestAlg3(t *testing.T) {
	sh := stdHandlerR4{
		// V: 2,
		ID0: string([]byte{0x4e, 0x00, 0x99, 0xe5, 0x36, 0x78, 0x93, 0x24,
			0xff, 0xd5, 0x82, 0xe4, 0xec, 0x0e, 0xa3, 0xb4}),
		Length: 128,
	}

	Oexp := []byte{0xE6, 0x00, 0xEC, 0xC2, 0x02, 0x88, 0xAD, 0x8B,
		0x0d, 0x64, 0xA9, 0x29, 0xC6, 0xA8, 0x3E, 0xE2, 0x51,
		0x76, 0x79, 0xAA, 0x02, 0x18, 0xBE, 0xCE, 0xEA, 0x8B, 0x79, 0x86,
		0x72, 0x6A, 0x8C, 0xDB}
	O, err := sh.alg3(3, []byte(""), []byte("test"))
	if err != nil {
		t.Errorf("crypt alg3 error %s", err)
		return
	}

	if string(O) != string(Oexp) {
		common.Log.Debug("   O (%d): % x", len(O), O)
		common.Log.Debug("Oexp (%d): % x", len(Oexp), Oexp)
		t.Errorf("alg3 -> key != expected")
	}
}

// Test algorithm 5 for computing dictionary's U (user password) value
// valid for R >= 3.
func TestAlg5(t *testing.T) {
	sh := stdHandlerR4{
		// V: 2,
		ID0: string([]byte{0x4e, 0x00, 0x99, 0xe5, 0x36, 0x78, 0x93, 0x24,
			0xff, 0xd5, 0x82, 0xe4, 0xec, 0x0e, 0xa3, 0xb4}),
		Length: 128,
	}
	d := &StdEncryptDict{
		R:               3,
		P:               0xfffff0c0,
		EncryptMetadata: true,
		O: []byte{0xE6, 0x00, 0xEC, 0xC2, 0x02, 0x88, 0xAD, 0x8B,
			0x5C, 0x72, 0x64, 0xA9, 0x5C, 0x29, 0xC6, 0xA8, 0x3E, 0xE2, 0x51,
			0x76, 0x79, 0xAA, 0x02, 0x18, 0xBE, 0xCE, 0xEA, 0x8B, 0x79, 0x86,
			0x72, 0x6A, 0x8C, 0xDB},
	}

	ekey := sh.alg2(d, []byte(""))
	U, err := sh.alg5(ekey, []byte(""))
	if err != nil {
		t.Errorf("Error %s", err)
		return
	}

	Uexp := []byte{0x59, 0x66, 0x38, 0x6c, 0x76, 0xfe, 0x95, 0x7d, 0x3d,
		0x0d, 0x14, 0x3d, 0x36, 0xfd, 0x01, 0x3d, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	if string(U[0:16]) != string(Uexp[0:16]) {
		common.Log.Info("   U (%d): % x", len(U), U)
		common.Log.Info("Uexp (%d): % x", len(Uexp), Uexp)
		t.Errorf("U != expected\n")
	}
}