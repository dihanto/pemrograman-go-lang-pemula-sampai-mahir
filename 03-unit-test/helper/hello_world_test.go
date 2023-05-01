package helper

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmakTable(b *testing.B) {
	benchmarks := []struct {
		name    string
		request string
	}{
		{
			name:    "Kurniawan",
			request: "Kurniawan",
		},
		{
			name:    "Kurniadi",
			request: "Kurniadi",
		},
		{
			name:    "KurniawanKurniadi",
			request: "Kurniawan Kurniadi",
		},
		{
			name:    "Asep",
			request: "Asep Kurniadi",
		},
	}
	for _, benchmark := range benchmarks {
		b.Run(benchmark.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				HelloWorld(benchmark.request)
			}
		})
	}
}

func BenchmarkSub(b *testing.B) {
	b.Run("Kurniawan", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			HelloWorld("Kurniawan")
		}
	})
	b.Run("Kurniadi", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			HelloWorld("Kurniadi")
		}
	})
}
func BenchmarkHelloWorld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HelloWorld("Kurniawan")
	}
}

func BenchmarkHelloWorldKurniadi(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HelloWorld("Kurniadi")
	}
}

func TestTableHelloWorld(t *testing.T) {
	tests := []struct {
		name     string
		request  string
		expected string
	}{
		{
			name:     "Kurniawan",
			request:  "Kurniawan",
			expected: "Hello Kurniawan",
		},
		{
			name:     "Kurniadi",
			request:  "Kurniadi",
			expected: "Hello Kurniadi",
		},
		{
			name:     "Asep",
			request:  "Asep",
			expected: "Hello Asep",
		},
		{
			name:     "Jamal",
			request:  "Jamal",
			expected: "Hello Jamal",
		},
		{
			name:     "Budi",
			request:  "Budi",
			expected: "Hello Budi",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := HelloWorld(test.request)
			require.Equal(t, test.expected, result)
		})

	}
}

func TestSubTest(t *testing.T) {
	t.Run("Kurniawan", func(t *testing.T) {
		result := HelloWorld("Kurniawan")
		require.Equal(t, "Hello Kurniawan", result, "Result Must be 'Hello Kurniawan'")
	})
	t.Run("Kurniadi", func(t *testing.T) {
		result := HelloWorld("Kurniadi")
		require.Equal(t, "Hello Kurniadi", result, "Result Must be 'Hello Kurniadi'")
	})
	t.Run("Asep", func(t *testing.T) {
		result := HelloWorld("Asep")
		require.Equal(t, "Hello Asep", result, "Result must be 'Hello Asep'")
	})
}
func TestMain(m *testing.M) {
	//before
	fmt.Println("Before unit test")

	m.Run()

	//after
	fmt.Println("After unit test")
}
func TestSkip(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("Can not run on Mac OS")
	}

	result := HelloWorld("Kurniawan")
	require.Equal(t, "Hello Kurniawan", result, "Result must be 'Hello Kurniawan'")
}
func TestHelloWorldRequire(t *testing.T) {
	result := HelloWorld("Kurniawan")
	require.Equal(t, "Hello Kurniawan", result, "Result must be 'Hello Kurniawan'")
	fmt.Println("TestHelloWorld with Require Done")
}

func TestHelloWorldAssert(t *testing.T) {
	result := HelloWorld("Kurniawan")
	assert.Equal(t, "Hello Kurniawan", result, "Result must be 'Hello Kurniawan'")
	fmt.Println("TestHelloWorld with Assert Done")
}
func TestHelloWorldKurniawan(t *testing.T) {
	result := HelloWorld("Kurniawan")

	if result != "Hello Kurniawan" {
		t.Error("Result must be 'Hello Kurniawan'")
	}

	fmt.Println("TestHelloWorldKurniawan Done")
}
func TestHelloWorldKhannedy(t *testing.T) {
	result := HelloWorld("Khannedy")

	if result != "Hello Khannedy" {
		// error
		t.Fatal("Result must be 'Hello Khannedy'")
	}

	fmt.Println("TestHelloWorldKhannedy Done")
}

func TestHelloWorldAsep(t *testing.T) {
	result := HelloWorld("Kurniawan")

	if result != "Hello Asep" {
		// error
		t.Fatal("Result must be 'Hello Asep'")
	}

	fmt.Println("TestHelloWorldAsep Done")
}
