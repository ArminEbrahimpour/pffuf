# pffuf
Parse fuzz faster u fool(ffuf)
## description
#### if you are using the [ffuf](https://github.com/ffuf/ffuf)  tool for fuzzing stuff you might use the json out put from it and well it is so confusing to work with a large json 
#### so in that case I wrote this tool 
#### which I called it pffuf to parse the json output form ffuf .

#### this is a basic tool which you can customize it easily and get parse different data out of the json 

## how to run ???
#### u just need to clone this repo or copy the code manually 
`git clone https://github.com/Bl00dBlu35/pffuf`
#### if youre using linux do these 
`go build pffuf.go `<br>
`sudo mv pffuf /usr/local/bin`<br> 
`// enter your password`
#### then you can run the tool where ever you want from terminal

## how does the program work
#### pffuf basically has one flag and it is `-f` which gets the path of the json file and make three directory based on http status code (302 for redirect, 403 for not allowed , 500 for server error endpoints)
#### and in each directory it writes the endpoints with that status code 

`pffuf -f ~/path/to/your/file.txt `
