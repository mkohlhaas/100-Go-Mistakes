package timer

import "testing"

func BenchmarkFoo1(b *testing.B) {
	expensiveSetup()

	for b.Loop() {
		functionUnderTest()
	}
}

func BenchmarkFoo2(b *testing.B) {
	for b.Loop() {
		b.StopTimer()
		expensiveSetup()
		b.StartTimer()
		functionUnderTest()
	}
}

func functionUnderTest() {
}

func expensiveSetup() {
}
