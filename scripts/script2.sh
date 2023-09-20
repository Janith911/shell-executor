# !/bin/bash

# out=$(nc -zV localhost 8081 | grep succeeded)
# if [ "$out" != "" ]
# then
#     echo SUCCESS
# else
#     echo FAILED
# fi

nc -z localhost 8080