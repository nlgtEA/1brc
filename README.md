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
|1| Naive approach, read data to map and process sequentially | 4:57.19 | - | - |

## 🚀 Bechmark Results
The overall results above is ran on my machine with many other running processes so it might vary.
This bechmark is ran using go testing bench, on a file with 1M rows, so it's supposed to be more stable and reliable.

```bash
make bench
```


| Attempt Number | Evaluate 1M | Evaluate 10M |
|----------------|---|---|
| 1 | 289ms |  |
