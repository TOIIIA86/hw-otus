package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkGetUsers(b *testing.B) {
	b.Helper()
	b.StopTimer()

	r, _ := zip.OpenReader("testdata/users.dat.zip")
	defer r.Close()

	data, _ := r.File[0].Open()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		getUsers(data)
	}
	b.StopTimer()
}

func BenchmarkCountDomains(b *testing.B) {
	b.Helper()
	b.StopTimer()

	r, _ := zip.OpenReader("testdata/users.dat.zip")
	defer r.Close()

	data, _ := r.File[0].Open()
	u, _ := getUsers(data)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		countDomains(u, "biz")
	}
	b.StopTimer()
}
