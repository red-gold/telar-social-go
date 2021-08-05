#!/bin/bash

# Colors                 
Black='\033[0;30m'        # Black
BRed='\033[1;31m'         # Red
BGreen='\033[1;32m'       # Green
BYellow='\033[1;33m'      # Yellow
BBlue='\033[1;34m'        # Blue
BPurple='\033[1;35m'      # Purple
BCyan='\033[1;36m'        # Cyan

opArg="${1}" 
arg="${2}"

echo Argument $arg1
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
echo -e "$BGreen Working directory: $DIR\n"
echo -e "$Black"

micros=(circles comments gallery posts user-rels vang votes)

# Update telar-core module
update_telar_core(){
for i in "${micros[@]}"; 
do 
echo -e "$BYellow"
microPath="$DIR/micros/$i"
echo "Updating [telar-core] module in $microPath"
cd "$microPath" && go get github.com/red-gold/telar-core@v0.1.16 && go mod tidy
echo "[telar-core] module updated for $i"
done
echo -e "$Black"
}


# Update ts-serverless module
update_ts_serverless(){
for i in "${micros[@]}"; 
do 
echo -e "$BYellow"
microPath="$DIR/micros/$i"
echo "Updating [ts-serverless] module in $microPath"
cd "$microPath" && go get github.com/red-gold/ts-serverless@v0.1.33 && go mod tidy
echo "[ts-serverless] module updated for $i"
done
echo -e "$Black"
}

case $opArg in

  telar-core)
    echo -n "Updating [telar-core] module ..."
    update_telar_core
    ;;
  ts-serverless)
    echo -n "Updating [ts-serverless] module"
    update_ts_serverless
    ;;

  *)
    echo -n "unknown"
    ;;
esac