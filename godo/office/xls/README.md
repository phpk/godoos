# xls

[![GoDoc](https://godoc.org/github.com/extrame/xls?status.svg)](https://godoc.org/github.com/extrame/xls)

Pure Golang xls library writen by [Rongshu Tech (chinese)](http://www.rongshu.tech), based on libxls. 

Thanks for contributions from Tamás Gulácsi @tgulacsi, @flyin9.

# Basic Usage

* Use **Open** function for open file
* Use **OpenWithCloser** function for open file and use the return value closer for close file
* Use **OpenReader** function for open xls from a reader, you should close related file in your own code

* Follow the example in GoDoc