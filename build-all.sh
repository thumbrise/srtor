BIN_FILE_NAME_PREFIX=$1
PROJECT_DIR=$2
PLATFORMS=$(go tool dist list)
for PLATFORM in $PLATFORMS; do
        GOOS=${PLATFORM%/*}
        GOARCH=${PLATFORM#*/}
        FILEPATH="$PROJECT_DIR/artifacts/${GOOS}-${GOARCH}"
        #echo $FILEPATH
        mkdir -p $FILEPATH
        BIN_FILE_NAME="$FILEPATH/${BIN_FILE_NAME_PREFIX}"
        #echo $BIN_FILE_NAME
        if [[ "${GOOS}" == "windows" ]]; then BIN_FILE_NAME="${BIN_FILE_NAME}.exe"; fi
        CMD="GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BIN_FILE_NAME}"
        #echo $CMD
        echo "${CMD}"
        eval $CMD || FAILURES="${FAILURES} ${PLATFORM}"
done