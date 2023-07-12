# edge-detection
A better edge detection algorithm written in go. 
Can be used easily from the ClI with one simple command:
```
./edge-detection [FILEPATH]
```
If you need to build it from source, simply do:
```
cd path/to/this/program
go build
```
## How it works
Uses the [Sobel Operator](https://en.wikipedia.org/wiki/Sobel_operator) to determine edges. 
