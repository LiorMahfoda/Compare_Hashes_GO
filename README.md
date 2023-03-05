# Compare_Hashes in Golang

In this script, I'm searching for hash values of files from a given csv file(Blacklist)\
and compers it to all files(recursively) from a given path\
If there is a match - the program askes the user wheather to delete the file/s or not\
There is a text file wit the output when the program terminated.\

In src.go file: run cmd with the following flags\
-p : path of search\
-n : path of the Blacklist.csv file(can be located anywhere in os file system)\
\
To run in cmd :\
go run src --path "Path to search" --name "Path to Blacklist.csv file"\
for example:\
go run src.go -p "C:\Users\Lior Mahfoda\Downloads\go_testing" -n "C:\Users\Lior Mahfoda\Downloads\go_testing\Blacklist.csv"\

To make it simple - there is a yaml file(config.yaml) that we can specify the path and name flags\
instead of rewriting it to the cmd\
Just run the following: "go run src.go" after putting the right paths in the config.yaml file.
