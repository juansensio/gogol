# GOGOL

Game of Life in GO

- [x] basic implementation with arrays, for loops and if-else statements
- [x] print in terminal
- [ ] output video
    - [x] save images
    - [ ] generate video
- [x] use structs and methods / interfaces
- [ ] organize code in main with cli args and GOL logic
- [x] pad arrays to avoid if-else statements
- [x] decouple simulation and visualization
- [x] matrix-vector multiply
- [x] sparse matrix
- [x] naive concurrency / parallelism
- [x] matrix-vector multiply paralel
- [x] benchmark (compare simulation speed of different implementations at different world sizes with same seed)

## Results

- Padding improves naive
- Using structs / methods make it slower.
- sparse matrix is faster than naive
- parallelism makes it faster (matrix > naive)