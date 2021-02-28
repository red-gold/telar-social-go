#!/bin/bash

# Colors                 
Black='\033[0;30m'        # Black
BRed='\033[1;31m'         # Red
BGreen='\033[1;32m'       # Green
BYellow='\033[1;33m'      # Yellow
BBlue='\033[1;34m'        # Blue
BPurple='\033[1;35m'      # Purple
BCyan='\033[1;36m'        # Cyan

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
echo -e "$BGreen Working directory: $DIR\n"
echo -e "$Black"

handlers=($(find $DIR/micros -maxdepth 5 -name '*go.mod'))
replace_line="//replace\nreplace github.com/red-gold/ts-serverless => $DIR"
lik='replace github.com/Qolzam'
current_file=""
len=${#handlers[*]}  #determine length of array

# Update go package function
update_go(){
  
  awk '/\/\/replace/{n=2}; n {n--; next}; 1' $current_file > $current_file.tmp && mv $current_file.tmp $current_file
  echo -e "$BGreen Remove replace from $current_file"
  echo -e "$Black"
}

if [ "$1" == "u" ]; then
# Update go packages
for (( i=0; i<len; i++ ))
do
    file=${handlers[$i]/go.mod/}
    cd $file
    go get -u github.com/red-gold/ts-serverless && go mod vendor
    echo -e "$BGreen Update vendor $file"
    echo -e "$Black"
done

elif [ "$1" == "r" ]; then
# Remove replace syntax
for (( i=0; i<len; i++ ))
do
    current_file=${handlers[$i]}
    update_go &

done

elif [ "$1" == "s" ]; then

# Add replace syntax
for (( i=0; i<len; i++ ))
do
    file=${handlers[$i]}
    echo -e "\n$replace_line" >> $file
    echo -e "$BGreen Add replace to $file"
    echo -e "$Black"

done

else 

echo -e "$BRed No command found:\n[u: update handlers vendor]\n[r: remove replace from go.mod] \n[s: set replace for go.mod] \n"
echo -e "$Black"

fi