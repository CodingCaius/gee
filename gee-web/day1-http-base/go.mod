module example

go 1.20

require gee v0.0.0

//替换，使用本地的gee库代替远程库
replace gee => ./gee