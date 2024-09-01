# Image-Compare
Image comparison tool written in golang.

## Usage
Usage of image-compare:\
  -A string\
        Filepath/directory A\
  -B string\
        Filepath/directory B\
  -c string\
        Optional: Comparison options, [pixel,contrast,quad]\ Default: all\
  -o string\
        Optional: output directory\

Ex. ```compare.exe -A ./red.png -B ./blue.png -c pixel -o ./results```