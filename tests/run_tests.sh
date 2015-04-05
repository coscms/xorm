#!/bin/bash
if [ "$#" == "0" ];then
    tests=('sqlite3' 'mysql' 'mymysql' 'postgres' 'mssql')
else
	tests=($@)
fi

for i in "${tests[@]}"
do :
   cd $i
   go test -v
   ec=$?
   cd -
   if [ $ec != "0" ];then
   		exit $ec
   fi
done

