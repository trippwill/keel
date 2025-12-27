# Keel Benchmarks

## Environment
- go version go1.25.5 linux/amd64
- cpu: Intel(R) Core(TM) Ultra 5 228V

## Render Baseline
Command:
`go test ./keel/... -bench=BenchmarkRender -benchmem`

Result:
`BenchmarkRenderExampleSplit-8 8188 150928 ns/op 37957 B/op 729`

## Resolver Baseline
Command:
`go test ./keel/... -bench=BenchmarkResolve -benchmem`

Results:
`BenchmarkResolveExtentAt/n=3-8 29659454 45.14 ns/op 24 B/op 1
BenchmarkResolveExtentAt/n=32-8 4436094 282.3 ns/op 256 B/op 1`
`BenchmarkResolveExtentAt/n=8-8 12841800 95.21 ns/op 64 B/op 1`
`BenchmarkResolveExtentAt/n=32-8 4436094 282.3 ns/op 256 B/op 1`
`BenchmarkResolveExtentAt/n=128-8 992932 1114 ns/op 1024 B/op 1`
`BenchmarkResolveCachedExtents/n=3-8 18715756 76.41 ns/op 104 B/op 2
BenchmarkResolveCachedExtents/n=32-8 3736586 318.6 ns/op 1024 B/op 2`
`BenchmarkResolveCachedExtents/n=8-8 10190673 116.5 ns/op 256 B/op 2`
`BenchmarkResolveCachedExtents/n=32-8 3736586 318.6 ns/op 1024 B/op 2`
`BenchmarkResolveCachedExtents/n=128-8 1000000 1151 ns/op 4096 B/op 2`
