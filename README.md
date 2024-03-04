# 1BRC

1️⃣🐝🏎️ [The One Billion Row Challenge](https://github.com/gunnarmorling/1brc) -- A fun exploration of how quickly 1B rows from a text file can be aggregated. The challenge was primarily foces on Java but I decided to solve it in Golang!


The original solution with the tests I re-use for this repository:
> I wrote a detailed blog about my implementation approach, you can check it out [here](https://www.bytesizego.com/blog/one-billion-row-challenge-go).


The benchmarking results are run on my personal MacBook Pro 2020, 1.4 GHz Quad-Core Intel Core i5, 8 GB 2133 MHz LPDDR3.

The measurement file is about 13GB in size.

As a baseline, the original naive Java implementation took 3:55.68 to run on this machine.

## 🚀 Results

| Attempt Number | Approach | Execution Time (m:ss) | Diff | Commit |
|-----------------|---|---|---|--|
|0| Naive Implementation: Read the file line by line into a map with key is the city name, and value is an array of 4 int of min, sum, max, avg value respectively (the original temp is scaled by a factor of 10 to avoid the cumbersone from adding and rounding float numbers) | 7:11.47 | | [753528e8](https://github.com/nlgtEA/1brc/commit/753528e8ac928a9525c60cfc648d3f3329dd631b)|
|1| Get rid of slice append and use direct access instead | 6:54.97 | 16.5 ||


## 🚀 Bechmark Results
The overall results above is ran on my machine with many other running processes so it might vary.
This bechmark is ran using go testing bench, on a file with 1M rows, so it's supposed to be more stable and reliable.

```bash
go test -bench=. Evaluate -count=1 -cpu=4
```


| Attempt Number | Time (ns/op) | Diff |
|----------------|---|---|
|0| TBU | |
|0| 308996589 | |
