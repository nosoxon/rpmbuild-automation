#!/bin/bash

lib=(ld-musl-x86_64.so.1 libacl.so.1 libpopt.so.0 libz.so.1)

usr_lib=(
	libbrotlicommon.so.1 libbrotlidec.so.1 libbz2.so.1 libcap.so.2
	libcreaterepo_c.so.0 libcrypto.so.1.1 libcurl.so.4 libgcrypt.so.20
	libglib-2.0.so.0 libgpg-error.so.0 libintl.so.8 liblzma.so.5 libmagic.so.1
	libnghttp2.so.14 libpcre.so.1 librpmio.so.9 librpm.so.9 libsqlite3.so.0
	libssl.so.1.1 libxml2.so.2 libzstd.so.1
)

mkdir -p /slice/{bin,lib,usr/lib/rpm}

for so in "${lib[@]}"; do
	cp -v /lib/$so /slice/lib
done

for so in "${usr_lib[@]}"; do
	cp -v /usr/lib/$so /slice/lib
done

cp /bin/rm /slice/bin
cp /usr/bin/createrepo_c /slice/bin/createrepo
touch /slice/usr/lib/rpm/rpmrc
