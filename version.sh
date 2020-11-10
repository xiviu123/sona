#!/usr/bin/env bash
versionFile="version"
ver=0.1
[[ ! -f ${versionFile} ]] && echo ${ver} > ${versionFile} && echo ${ver} && exit
ver=`cat ${versionFile}`
newVer=`python -c "print($ver+0.1)"`
echo ${newVer} > ${versionFile} && echo ${newVer}
