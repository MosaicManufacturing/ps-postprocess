cd "./licenses" || exit
go build || exit
./licenses update
