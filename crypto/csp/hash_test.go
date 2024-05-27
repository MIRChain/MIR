package csp

import (
	"bytes"
	"encoding/asn1"
	"fmt"
	"testing"
)

func TestHash_Sum(t *testing.T) {
	buf := new(bytes.Buffer)
	for _, testStr := range []string{"", "some test string"} {
			for _, algo := range []asn1.ObjectIdentifier{GOST_R3411, GOST_R3411_12_256, GOST_R3411_12_512} {
			func() {
				h, err := NewHash(HashOptions{HashAlg: algo})
				if err != nil {
					t.Error(err)
					return
				}
				defer func() {
					if err := h.Close(); err != nil {
						t.Error(err)
					}
				}()
				fmt.Fprintf(h, "%s", testStr)
				fmt.Fprintf(buf, "%s %d %q %x\n", algo, h.Size()*8, testStr, h.Sum(nil))
			}()
		}
	}
}

func TestHash_HMAC_Sum(t *testing.T) {
	buf := new(bytes.Buffer)
	for _, algo := range []asn1.ObjectIdentifier{GOST_R3411, GOST_R3411_12_256, GOST_R3411_12_512} {
		for _, testKey := range []string{"", "1234", "some other key"} {
			for _, testStr := range []string{"", "The quick brown fox jumps over the lazy dog"} {
				func() {
					h, err := NewHMAC(algo, ([]byte)(testKey))
					if err != nil {
						t.Error(err)
						return
					}
					defer func() {
						if err := h.Close(); err != nil {
							t.Error(err)
						}
					}()
					fmt.Fprintf(buf, "%s %q %q %x\n", algo, testKey, testStr, h.Sum(([]byte)(testStr)))
				}()
			}
		}
	}
}