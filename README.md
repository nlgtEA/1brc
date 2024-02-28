# 1BRC

1️⃣🐝🏎️ [The One Billion Row Challenge](https://github.com/gunnarmorling/1brc) -- A fun exploration of how quickly 1B rows from a text file can be aggregated. The challenge was primarily foces on Java but I decided to solve it in Golang!


The original solution with the tests I re-use for this repository:
> I wrote a detailed blog about my implementation approach, you can check it out [here](https://www.bytesizego.com/blog/one-billion-row-challenge-go).


The benchmarking results are run on my personal MacBook Pro 2020, 1.4 GHz Quad-Core Intel Core i5, 8 GB 2133 MHz LPDDR3.
The measurement file is about 13GB in size.
As a baseline, the original naive Java implementation took 3:55.68 to run on this machine.

## 🚀 Results

| Attempt Number | Approach | Execution Time | Diff | Commit |
|-----------------|---|---|---|--|
|0| Naive Implementation: Read the file line by line into a map with key is the city name, and value is an array of 4 int of min, sum, max, avg value respectively (the original temp is scaled by a factor of 10 to avoid the cumbersone from adding and rounding float numbers) | 10:01.12 | | [753528e](https://github.com/nlgtEA/2brc/commit/753528e8ac928a9525c60cfc648d3f3329dd631b)|
