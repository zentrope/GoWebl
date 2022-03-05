#!/bin/zsh
# Shell Script to create all relevant icons from an high resolution artwork
# https://gist.github.com/pd95/f8afdeb6139d96019bd60b11a70d6367

if [ "x$1" != "x" -a -f "$1" ] ; then
    INPUT=$1
else 
    INPUT="Artwork.png"
fi

if [ ! -f "$INPUT" ]; then
    echo "Input file not specified"
    exit 1
fi
mkdir -p AppIcons

while IFS= read -r line ; do
  FILENAME=${line% *}
  SIZE=${line#* }
  if [ "x$FILENAME" != "x" ] ; then
    COMMAND="sips -s format png \"$INPUT\" --resampleHeightWidth $SIZE $SIZE --out \"AppIcons/$FILENAME\""
    echo "$COMMAND"
    eval "$COMMAND"
  else
    echo ""
  fi
done <<EOF
AppIcon~mac-512@2x.png 1024
AppIcon~mac-512@1x.png 512
AppIcon~mac-256@2x.png 512
AppIcon~mac-256@1x.png 256
AppIcon~mac-128@2x.png 256
AppIcon~mac-128@1x.png 128
AppIcon~mac-32@2x.png 64
AppIcon~mac-32@1x.png 32
AppIcon~mac-16@2x.png 32
AppIcon~mac-16@1x.png 16
EOF

# AppIcon~mac-29@1x.png 29
# AppIcon~mac-29@2x.png 58
# AppIcon~mac-29@3x.png 87
# AppIcon~mac-40@2x.png 80
# AppIcon~mac-40@3x.png 120
# AppIcon~mac-57@1x.png 57
# AppIcon~mac-57@2x.png 114
# AppIcon~mac-60@2x.png 120
# AppIcon~mac-60@3x.png 180
# 
# AppIcon~iPad-20@1x.png 20
# AppIcon~iPad-20@2x.png 40
# AppIcon~iPad-29@1x.png 29
# AppIcon~iPad-29@2x.png 58
# AppIcon~iPad-40@1x.png 40
# AppIcon~iPad-40@2x.png 80
# AppIcon~iPad-76@1x.png 76
# AppIcon~iPad-76@2x.png 152
# AppIcon~iPad-83.5@2x.png 167
# 
# AppIcon~iOS-Marketing-1024@1x.png 1024
