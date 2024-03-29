# 1BRC

1️⃣🐝🏎️ [The One Billion Row Challenge](https://github.com/gunnarmorling/1brc) -- A fun exploration of how quickly 1B rows from a text file can be aggregated. The challenge was primarily foces on Java but I decided to solve it in Golang!


The original solution with the tests I re-use for this repository:
> I wrote a detailed blog about my implementation approach, you can check it out [here](https://www.bytesizego.com/blog/one-billion-row-challenge-go).


The benchmarking results are run on my personal MacBook Pro 2020, 1.4 GHz Quad-Core Intel Core i5, 8 GB 2133 MHz LPDDR3.

The measurement file is about 13GB in size.

As a baseline, the original naive Java implementation took 3:55.68 to run on this machine.

## 🚀 Results

| Attempt Number | Approach | Execution Time (m:ss) | Note | Commit |
|-----------------|---|---|---|--|
|1| Naive approach, read data to map and process sequentially | 4:57.19 | - | - |
|2| Using readSlice and custom parseInt function  | 4:18.95 | - | - |
|3| Read file by chunks of 20MB  | 4:07.94 | - | - |
|3.1| Read file by chunks of 100MB  | 4:14.20 | - | - |
|3.1| Read file by chunks of 50MB  | 4:09.89 | - | - |
|4| Using concurrency with 6 worker in total including main thread  | 1:27.80 | The chunk size of 20MB and 6 gouroutines seems to be the sweet spot. Other set up are drastically slower, usually take around 3m50s. And not sure why it's not faster in the case of 1M row file. | - |
|5| Replace bytes.Split with finding ; manually  | 52.513 | - | - |
|6| Reduce hash map read, and find new optimal chunk size  | 44.495 | - | - |

## 🚀 Bechmark Results
The overall results above is ran on my machine with many other running processes so it might vary.
This bechmark is ran using go testing bench, on a file with 1M rows, so it's supposed to be more stable and reliable.

```bash
make bench
```


| Attempt Number | Evaluate 1M | Evaluate 10M |
|----------------|---|---|
| 1 | 289ms |  |
| 2 | 245ms |  |
| 3 | 254ms |  |
| 4 | 241ms |  |
| 5 | 175ms |  |
| 6 | 117ms |  |
