package l_test

import (
	"os"
	"testing"

	"example.com/l"
)

func TestSimpleErrorToFile(t *testing.T) {
	_, err := os.OpenFile("test.txt", os.O_RDONLY, 0)

	for i := 0; i < 100_000; i++ {
		l.ErrorTrace(err)
	}
}

func BenchmarkSimpleErrorStack(b *testing.B) {
	_, err := os.OpenFile("test.txt", os.O_RDONLY, 0)
	os.Remove("logs/zero.log")

	for i := 0; i < b.N; i++ {
		for i := 0; i < 100_000; i++ {
			l.ErrorTrace(err)
		}
	}
}

func BenchmarkSimpleErrorNoStack(b *testing.B) {
	_, err := os.OpenFile("test.txt", os.O_RDONLY, 0)
	os.Remove("logs/zero.log")

	for i := 0; i < b.N; i++ {
		for i := 0; i < 100_000; i++ {
			l.Error(err)
		}
	}
}
