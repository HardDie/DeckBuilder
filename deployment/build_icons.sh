#!/bin/bash

mkdir icons
cp 512.png icons/icon_512x512.png
convert -resize 256x 512.png icons/icon_256x256.png
convert -resize 128x 512.png icons/icon_128x128.png
convert -resize 32x 512.png icons/icon_32x32.png
convert -resize 16x 512.png icons/icon_16x16.png

#libicns
png2icns icon.icns icons/*
