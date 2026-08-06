package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	hello1Hash := hello1.HashKey()
	hello2Hash := hello2.HashKey()
	diff1Hash := diff1.HashKey()
	diff2Hash := diff2.HashKey()

	if hello1Hash != hello2Hash {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1Hash != diff2Hash {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1Hash == diff1Hash {
		t.Errorf("strings with different content have same hash keys")
	}
}
