# Image-Compare
Image comparison tool written in golang.

## Usage
### Compare
```
Usage of compare:
  -A string
        Filepath/directory A.
  -B string
        Filepath/directory B.
  -c string
        Optional: Comparison options, [pixel,contrast,quad,ssim,mse]. (default "all")
  -o string
        Optional: output directory. 
```
Ex. ```compare.exe -A ./red.png -B ./blue.png -c pixel -o ./results```
### Filter
```
Usage of filter:
  -c string
        Optional: Comparison options, [pixel,contrast,quad,ssim,mse]. (default "all")
  -d string
        Optional: Path to directory to filter.
  -i float
        Optional: Index threshold. (default 1)
  -n int
        Optional: Num failed points.
```
Ex. ```filter.exe -d ./result/ -i 0.99 -c ssim```