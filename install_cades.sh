#!/usr/bin/env bash

archive='cades_linux_amd64.tar.gz'
tmp_dir='cadesinstall'

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
echo "Installing Cades plugin"
cd /${tmp_dir}/cades*/
rpm -Uvh cprocsp-pki-cades*.rpm cprocsp-pki-plugin*.rpm
cd -
echo "Deleting tmp dir"
rm -r ${tmp_dir} 

