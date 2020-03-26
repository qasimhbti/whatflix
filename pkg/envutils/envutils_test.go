package envutils

import "testing"

func TestCheckEnvironment(t *testing.T) {
	for _, env := range []string{Testing, Development, Staging, Production} {
		err := Check(env)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCheckEnvironmentError(t *testing.T) {
	err := Check("invalid")
	if err == nil {
		t.Fatal("no error")
	}
}
