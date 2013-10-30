/*
Copyright 2013 Tristan Wietsma

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package jack

import (
	"fmt"
	"testing"
)

//////////////
// Examples //
//////////////

func ExampleGet() {
	m := Message{GET, []byte("key123"), nil}
	fmt.Println(m.Bytes())
	// Output:
	// [103 30 107 101 121 49 50 51 0]
}

func ExampleParseGet() {
	b := []byte{103, 30, 107, 101, 121, 49, 50, 51, 0}
	m, _ := Parse(b)
	fmt.Println(*m)
	// Output:
	// {103 [107 101 121 49 50 51] []}
}

func ExampleSet() {
	m := Message{SET, []byte("key123"), []byte("val567")}
	fmt.Println(m.Bytes())
	// Output:
	// [115 30 107 101 121 49 50 51 30 118 97 108 53 54 55 0]
}

func ExampleParseSet() {
	b := []byte{115, 30, 107, 101, 121, 49, 50, 51, 30, 118, 97, 108, 53, 54, 55, 0}
	m, _ := Parse(b)
	fmt.Println(*m)
	// Output:
	// {115 [107 101 121 49 50 51] [118 97 108 53 54 55]}
}

func ExampleDel() {
	m := Message{DEL, []byte("key123"), nil}
	fmt.Println(m.Bytes())
	// Output:
	// [100 30 107 101 121 49 50 51 0]
}

func ExampleParseDel() {
	b := []byte{100, 30, 107, 101, 121, 49, 50, 51, 0}
	m, _ := Parse(b)
	fmt.Println(*m)
	// Output:
	// {100 [107 101 121 49 50 51] []}
}

func ExamplePub() {
	m := Message{PUB, []byte("key123"), []byte("val567")}
	fmt.Println(m.Bytes())
	// Output:
	// [80 30 107 101 121 49 50 51 30 118 97 108 53 54 55 0]
}

func ExampleParsePub() {
	b := []byte{80, 30, 107, 101, 121, 49, 50, 51, 30, 118, 97, 108, 53, 54, 55, 0}
	m, _ := Parse(b)
	fmt.Println(*m)
	// Output:
	// {80 [107 101 121 49 50 51] [118 97 108 53 54 55]}
}

func ExampleSub() {
	m := Message{SUB, []byte("key123"), nil}
	fmt.Println(m.Bytes())
	// Output:
	// [83 30 107 101 121 49 50 51 0]
}

func ExampleParseSub() {
	b := []byte{83, 30, 107, 101, 121, 49, 50, 51, 0}
	m, _ := Parse(b)
	fmt.Println(*m)
	// Output:
	// {83 [107 101 121 49 50 51] []}
}

////////////////
// Benchmarks //
////////////////

func BenchmarkParseGet(b *testing.B) {
	bt := []byte{103, 30, 107, 101, 121, 49, 50, 51, 0}
	for i := 0; i < b.N; i++ {
		_, _ = Parse(bt)
	}
}

func BenchmarkParseSet(b *testing.B) {
	bt := []byte{115, 30, 107, 101, 121, 49, 50, 51, 30, 118, 97, 108, 53, 54, 55, 0}
	for i := 0; i < b.N; i++ {
		_, _ = Parse(bt)
	}
}

func BenchmarkParseDel(b *testing.B) {
	bt := []byte{100, 30, 107, 101, 121, 49, 50, 51, 0}
	for i := 0; i < b.N; i++ {
		_, _ = Parse(bt)
	}
}

func BenchmarkParsePub(b *testing.B) {
	bt := []byte{80, 30, 107, 101, 121, 49, 50, 51, 30, 118, 97, 108, 53, 54, 55, 0}
	for i := 0; i < b.N; i++ {
		_, _ = Parse(bt)
	}
}

func BenchmarkParseSub(b *testing.B) {
	bt := []byte{83, 30, 107, 101, 121, 49, 50, 51, 0}
	for i := 0; i < b.N; i++ {
		_, _ = Parse(bt)
	}
}

///////////
// Tests //
///////////

func TestMissingTerminalByte(t *testing.T) {
	b := []byte{83, 30, 107, 101, 121, 49, 50, 51}
	_, err := Parse(b)
	if err.Error() != "ProtocolError: Message is missing EOM byte." {
		t.Errorf("Missing an error.")
	}
}

func TestInvalidCommand(t *testing.T) {
	b := []byte{42, 30, 107, 101, 121, 49, 50, 51, 0}
	_, err := Parse(b)
	if err.Error() != "ProtocolError: Message invokes a nonexistent command." {
		t.Errorf("Missing an error.")
	}
}

func TestMissingKey(t *testing.T) {
	b := []byte{83, 0}
	_, err := Parse(b)
	if err.Error() != "ProtocolError: Message is missing key." {
		t.Errorf("Missing an error.")
	}
}

func TestNonsense(t *testing.T) {
	b := []byte{83, 30, 107, 101, 121, 49, 50, 51, 30, 51, 0}
	_, err := Parse(b)
	if err.Error() != "ProtocolError: Message is nonsense." {
		t.Errorf("Missing an error.")
	}
}
