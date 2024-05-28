#!/usr/bin/env bash

archive='linux-amd64_rpm.tgz'
tmp_dir='cspinstall'

if [[ ! -f ${archive} ]]
then
	echo "No CSP archive ${archive}!"
	exit 1
fi


function unpack {
	mkdir -p ${tmp_dir}
	echo "Unpacking ${archive}"
	tar xzf ${archive} -C ${tmp_dir}
}

unpack
echo "Installing CSP"
./${tmp_dir}/linux*/install.sh
./${tmp_dir}/linux*/install.sh kc1 lsb-cprocsp-devel cprocsp-stunnel-msspi
rpm -Uvh ./${tmp_dir}/linux*/cprocsp-rdr-gui-gtk*.rpm
echo "Deleting tmp dir"
rm -r ${tmp_dir} 
