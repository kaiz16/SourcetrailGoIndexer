ABSDIR=`pwd`/test
if [ "$1" != "" ]
then
	ABSDIR=$1
fi

echo "Indexing" $ABSDIR "..."

export LD_LIBRARY_PATH=./sourcetrail:$LD_LIBRARY_PATH
# export LD_DEBUG=help
go run ./src/ -pkgPath=$ABSDIR

sourcetrail $ABSDIR/cg.srctrlprj