/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

// Test the PDF crypt support.

package core

import (
	"testing"

	"github.com/gnaoh1379/unipdf/common"
	"github.com/gnaoh1379/unipdf/core/security"
)

func init() {
	common.SetLogger(common.ConsoleLogger{})
}

// Test decrypting. Example with V=2, R=3, using standard algorithm.
func TestDecryption1(t *testing.T) {
	crypter := PdfCrypt{
		encrypt: encryptDict{
			V:      2,
			Length: 128,
		},
		encryptStd: security.StdEncryptDict{
			R:               3,
			P:               0xfffff0c0,
			EncryptMetadata: true,
			O: []byte{0xE6, 0x00, 0xEC, 0xC2, 0x02, 0x88, 0xAD, 0x8B,
				0x0d, 0x64, 0xA9, 0x29, 0xC6, 0xA8, 0x3E, 0xE2, 0x51,
				0x76, 0x79, 0xAA, 0x02, 0x18, 0xBE, 0xCE, 0xEA, 0x8B, 0x79, 0x86,
				0x72, 0x6A, 0x8C, 0xDB},
			U: []byte{0xED, 0x5B, 0xA7, 0x76, 0xFD, 0xD8, 0xE3, 0x89,
				0x4F, 0x54, 0x05, 0xC1, 0x3B, 0xFD, 0x86, 0xCF, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00},
		},
		id0: string([]byte{0x5f, 0x91, 0xff, 0xf2, 0x00, 0x88, 0x13,
			0x5f, 0x30, 0x24, 0xd1, 0x0f, 0x28, 0x31, 0xc6, 0xfa}),
		// Default algorithm is V2 (RC4).
		cryptFilters:     newCryptFiltersV2(128),
		decryptedObjects: make(map[PdfObject]bool),
	}

	streamData := []byte{0xBC, 0x89, 0x86, 0x8B, 0x3E, 0xCF, 0x24, 0x1C,
		0xC4, 0x88, 0xF3, 0x60, 0x74, 0x8A, 0x22, 0xE3, 0xAD, 0xF4, 0x48,
		0x8E, 0x20, 0x94, 0x06, 0x4B, 0x4B, 0xB5, 0x3E, 0x93, 0x89, 0x4E,
		0x32, 0x38, 0xB4, 0xF6, 0x05, 0x3C, 0x5D, 0x0C, 0x12, 0xE4, 0xEB,
		0x9B, 0x8D, 0x26, 0x32, 0x7B, 0x09, 0x97, 0xA1, 0xC5, 0x98, 0xF6,
		0xE7, 0x1C, 0x3B}

	// Plain text stream (hello world).
	exp := []byte{0x20, 0x20, 0x42, 0x54, 0x0A, 0x20, 0x20, 0x20, 0x20,
		0x2F, 0x46, 0x31, 0x20, 0x31, 0x38, 0x20, 0x54, 0x66, 0x0A, 0x20,
		0x20, 0x20, 0x20, 0x30, 0x20, 0x30, 0x20, 0x54, 0x64, 0x0A, 0x20,
		0x20, 0x20, 0x20, 0x28, 0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x20, 0x57,
		0x6F, 0x72, 0x6C, 0x64, 0x29, 0x20, 0x54, 0x6A, 0x0A, 0x20, 0x20,
		0x45, 0x54}
	rawText := "2 0 obj\n<< /Length 55 >>\nstream\n" + string(streamData) + "\nendstream\n"

	parser := PdfParser{}
	parser.xrefs.ObjectMap = make(map[int]XrefObject)
	parser.objstms = make(objectStreams)
	parser.rs, parser.reader, parser.fileSize = makeReaderForText(rawText)
	parser.crypter = &crypter

	obj, err := parser.ParseIndirectObject()
	if err != nil {
		t.Errorf("Error parsing object")
		return
	}

	so, ok := obj.(*PdfObjectStream)
	if !ok {
		t.Errorf("Should be stream (is %q)", obj)
		return
	}

	authenticated, err := parser.Decrypt([]byte(""))
	if err != nil {
		t.Errorf("Error authenticating")
		return
	}
	if !authenticated {
		t.Errorf("Failed to authenticate")
		return
	}

	parser.crypter.Decrypt(so, 0, 0)
	if string(so.Stream) != string(exp) {
		t.Errorf("Stream content wrong")
		return
	}
}
