#!/bin/sh

# target app.
if [ $# != 3 ] ; then
    echo "Usage: $0 APPBIN APPARGS BINPATH"
    exit 1
fi

# variables.
APPBIN=$1
APPARGS=$2
BINPATH=$3

# set the sed cmd according to OS
SEDCMD=sed
if [ `uname` == "Darwin" ]; then
  SEDCMD=gsed
fi
echo $SEDCMD

# set the readlink cmd according to OS
READLINKCMD=readlink
if [ `uname` == "Darwin" ]; then
  READLINKCMD=greadlink
fi
echo ${READLINKCMD}

# escape character ‘/’ in app args.
APPARGS=$(echo "$APPARGS" | ${SEDCMD} 's/\//\\\//g')

# tools.
TOOLS=$(dirname $(${READLINKCMD} -f "$0"))

# generate app daemon-control tool.
cp -rf ${TOOLS}/daemon-control.sh ${BINPATH}/$1.sh

${SEDCMD} -i "s/.*APPBIN=.*/APPBIN=\"${APPBIN}\"/" ${BINPATH}/$1.sh
${SEDCMD} -i "s/.*APPARGS=.*/APPARGS=\"${APPARGS}\"/" ${BINPATH}/$1.sh
${SEDCMD} -i "s/.*BINPATH=.*/BINPATH=\".\"/" ${BINPATH}/$1.sh
