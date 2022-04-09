// Copyright 2022 The Serenity Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	program := Program(`>++++++++[-<+++++++++>]<.>>+>-[+]++>++>+++[>[->+++<<+++>]<<]>-----.>-> Comments can be added
  +++..+++.>-.<<+[>[+>+]>>]<--------------.>>.+++.------.--------.>+.>+.`)
	output := program.Execute()
	if strings.TrimSpace(output.String()) != "Hello World!" {
		t.Fatalf("incorrect output: '%s'", output.String())
	}
}
