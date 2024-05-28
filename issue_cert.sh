#!/bin/bash
echo "o" | /opt/cprocsp/bin/amd64/certmgr -inst -store uroot  -silent -all  -file /cacerts.p7b

/opt/cprocsp/bin/amd64/csptest -keys -enum -verifyco -deletek -pattern "\\.\HDIMAGE"
/opt/cprocsp/bin/amd64/certmgr -del -store uMy -all

rm $DATADIR/certs/*.cert

sh -c "echo 1 | ./req_user.sh"
sh -c "echo 2 | ./req_user.sh"
sh -c "echo 3 | ./req_user.sh"