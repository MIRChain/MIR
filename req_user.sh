#!/bin/bash
/opt/cprocsp/sbin/amd64/cpconfig -hardware rndm -del bio_tui
/opt/cprocsp/sbin/amd64/cpconfig -hardware rndm -add cpsd -name 'cpsd rng' -level 5
cp kis_1 /var/opt/cprocsp/dsrf/db1/kis_1
chmod +w /var/opt/cprocsp/dsrf/db1/kis_1

/opt/cprocsp/sbin/amd64/cpconfig -hardware rndm -configure cpsd -add string /db1/kis_1 /var/opt/cprocsp/dsrf/db1/kis_1

function register () {
        output=$(/opt/cprocsp/bin/amd64/cryptcp -creatuser -field "CN=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 10 | head -n 1)" -CPCA20 "https://testca2012.cryptopro.ru/ui/879a7919-656d-4690-82c2-aa0d008c75f8" | grep "Token\|Password")
        user=$(echo "$output" | grep "Token" | awk '{print $3}')
        pass=$(echo "$output" | grep "Password"| awk '{print $2}')
        creds="$user:$pass"
        echo $creds
}

{
if [ -z "$1" ]; then
    echo 'Выберите тип сертификата:'
    echo '1)  Signer'
    echo '2)  [TLS Server]'
    echo '3)  [TLS Client]'
    read TYPE
else
    TYPE=$1
fi

case $TYPE in
1)  TMPL=sig;;
2)  TMPL=tls-srv;;
3)  TMPL=tls-clt;;
esac

CACRED=$(register $DNS)
CAURL="https://${CACRED}@testca2012.cryptopro.ru/ui/"

CATOKEN=`echo $CAURL | sed -E 's#[^:]*://([^:]*):.*#\1#'`
CAPASSWORD=`echo $CAURL | sed -E 's#[^:]*://[^:]*:([^@]*)@.*#\1#'`
CAURL=`echo $CAURL | sed -E 's#([^:]*://)[^@]*@(.*)$#\1\2#'`

RDN='E=testing@mirchain.org'
RQID=`mktemp`
RES=`mktemp`

ABBR=`echo $TMPL | grep -o '\\-..' | tr -d '\n-'`
ABBR="${ABBR:0:6}"
ABBR="${ABBR,,}"
if [ ! -z "$2" ]; then
    NUM="$2"
else
    NUM=`/opt/cprocsp/bin/amd64/csptestf -keyset -verifyc -enum | grep -E "^[0-9]{2}$ABBR\$" | wc -l`
fi
echo 
CONT=`printf '%02d%s' "$NUM" "$ABBR"`

case "$TMPL" in
    sig)
    AT='-sg' ;;
    *)
        AT='-both' ;;

esac
#AT='-both' 
echo /opt/cprocsp/bin/amd64/cryptcp -creatcert -rdn $RDN \
    -provtype 80 "$AT" -ku -hashalg '1.2.643.7.1.1.2.2' \
    -cpca20 "$CAURL" -token "$CATOKEN" -tpassword "$CAPASSWORD" \
    -tmpl "$TMPL" -fileid "$RQID" -cont "$CONT"

echo -e '\n' | /opt/cprocsp/bin/amd64/cryptcp -creatcert -rdn $RDN \
    -exprt -provtype 80 "$AT" -ku -hashalg '1.2.643.7.1.1.2.2' \
    -cpca20 "$CAURL" -token "$CATOKEN" -tpassword "$CAPASSWORD" \
    -tmpl "$TMPL" -fileid "$RQID" -cont "\\\\.\\HDIMAGE\\$CONT"

echo "Waiting for certificate request $(<$RQID) to be processed"

while : ; do
    echo /opt/cprocsp/bin/amd64/cryptcp -pendcert -FileID "$RQID" -cont "$CONT" \
        -CPCA20 "$CAURL" -token "$CATOKEN" -tpassword "$CAPASSWORD" > "$RES" 2>&1
    /opt/cprocsp/bin/amd64/cryptcp -pendcert -FileID "$RQID" -cont "$CONT" \
        -CPCA20 "$CAURL" -token "$CATOKEN" -tpassword "$CAPASSWORD" > "$RES" 2>&1
    cat "$RES"
    egrep -q "installed|установлен" "$RES" && break
    echo "Waiting for certificate request $(<$RQID) to be processed"
    sleep 1
done

rm "$RQID" "$RES"
echo "/opt/cprocsp/bin/amd64/certmgr -list | grep $CONT -B6 | grep SubjKeyID  | awk '{print $3}'"
KEYID=`/opt/cprocsp/bin/amd64/certmgr -list | grep $CONT -B6 | grep SubjKeyID  | awk '{print $3}'`
echo "/opt/cprocsp/bin/amd64/certmgr -export -pfx   -keyid $KEYID   -provtype 80  -dest  $DATADIR/$TMPL.pfx  -pin $TMPL"
/opt/cprocsp/bin/amd64/certmgr -export -pfx   -keyid $KEYID   -provtype 80  -dest  $DATADIR/$TMPL.pfx  -pin $TMPL
echo "/opt/cprocsp/bin/amd64/certmgr -export -keyid $KEYID   -provtype 80  -dest  $DATADIR/$TMPL.cer "
/opt/cprocsp/bin/amd64/certmgr -export -keyid $KEYID   -provtype 80  -dest  $DATADIR/$TMPL.cer 
#/opt/cprocsp/bin/amd64/csptestf  -passwd -pass "" -change $TMPL  -provtype 80  -container  $CONT
#/opt/cprocsp/bin/amd64/csptestf -ipsec -reg -mycert $DATADIR/$TMPL.cer -autocont -savepin -passw $TMPL
echo "/opt/cprocsp/bin/amd64/csptestf  -passwd -pass "" -change $TMPL -provtype 80  -container  $(/opt/cprocsp/bin/amd64/certmgr -list | grep $KEYID -A6 | grep Container | awk '{print $3}')"
/opt/cprocsp/bin/amd64/csptestf  -passwd -pass "" -change $TMPL -provtype 80  -container  $(/opt/cprocsp/bin/amd64/certmgr -list | grep $KEYID -A6 | grep Container | awk '{print $3}')
echo "/opt/cprocsp/sbin/amd64/cpconfig -ini "\\LOCAL\\KeyDevices\\passwords\\HDIMAGE\\"$(/opt/cprocsp/bin/amd64/certmgr -list | grep $KEYID -A6 | grep Container | awk '{print $3}' | awk -F '\' '{print $3}') -add string passwd $TMPL "
/opt/cprocsp/sbin/amd64/cpconfig -ini "\\LOCAL\\KeyDevices\\passwords\\HDIMAGE\\"$(/opt/cprocsp/bin/amd64/certmgr -list | grep $KEYID -A6 | grep Container | awk '{print $3}' | awk -F '\' '{print $3}') -add string passwd $TMPL 
echo "/opt/cprocsp/bin/amd64/csptestf -passwd -showsaved -container "\\\\.\\HDIMAGE\\"$(/opt/cprocsp/bin/amd64/certmgr -list | grep $KEYID -A6 | grep Container | awk '{print $3}' | awk -F '\' '{print $3}' | awk -F '.' '{print $1}') "
/opt/cprocsp/bin/amd64/csptestf -passwd -showsaved -container "\\\\.\\HDIMAGE\\"$(/opt/cprocsp/bin/amd64/certmgr -list | grep $KEYID -A6 | grep Container | awk '{print $3}' | awk -F '\' '{print $3}' | awk -F '.' '{print $1}') 
echo $KEYID > $DATADIR/$TMPL

} 
