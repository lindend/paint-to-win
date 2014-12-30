#options...
#-o output dir
#-d docker repo to publish to
#-u docker repo user name
#-p docker repo password

SolutionRoot=`readlink -f ../../../`
ProjectName="wordlist"
OutPath="$SolutionRoot/bin"

while getopts o: opt; do
	case $opt in
	o)
		OutPath = $OPTARG
		;;
	esac
done
shift $((OPTIND-1))

OutFullPath="$OutPath/$ProjectName"

GOPATH="$SolutionRoot"
go build -o $OutFullPath

exit $?