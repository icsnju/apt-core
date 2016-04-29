#!/usr/bin/env bash

# Fail on error, verbose output
set -exo pipefail

# Build project
#ndk-build 1>&2
id=$2
# Figure out which ABI and SDK the device has
abi=$(adb -s $id shell getprop ro.product.cpu.abi | tr -d '\r')
sdk=$(adb -s $id shell getprop ro.build.version.sdk | tr -d '\r')
rel=$(adb -s $id shell getprop ro.build.version.release | tr -d '\r')

# PIE is only supported since SDK 16
if (($sdk >= 16)); then
  bin=minicap
else
  bin=minicap-nopie
fi

args=
if [ "$1" = "autosize" ]; then
  set +o pipefail
  #size=$(adb -s $id shell dumpsys window | grep -Eo 'init=\d+x\d+' | head -1 | cut -d= -f 2)
  size=$(adb -s $id shell dumpsys window | grep 'init' | cut -d= -f 2 | cut -d' ' -f 1)

  if [ "$size" = "" ]; then
    #w=$(adb -s $id shell dumpsys window | grep -Eo 'DisplayWidth=\d+' | head -1 | cut -d= -f 2)
	w=$(adb -s $id shell dumpsys display | grep mBaseDisplayInfo | cut -d, -f 2 | cut -d' ' -f 3)
    #h=$(adb -s $id shell dumpsys window | grep -Eo 'DisplayHeight=\d+' | head -1 | cut -d= -f 2)
	h=$(adb -s $id shell dumpsys display | grep mBaseDisplayInfo | cut -d, -f 2 | cut -d' ' -f 5)
    size="${w}x${h}"
  fi
#  zoom=$size
#  if [$h -gt $w];then
#	w2=$((${w}/(${h}/400)))
#	zoom="${w2}x400"
#  else
#	h2=$((${h}/(${w}/400)))
#	zoom="400x${h2}"
#  fi
	
  args="-P $size@500x500/0"
  set -o pipefail
  shift
fi

# Create a directory for our resources
dir=/data/local/tmp/minicap-devel
adb -s $id shell "mkdir $dir 2>/dev/null"

# Upload the binary
adb -s $id push minicap/libs/$abi/$bin $dir

# Upload the shared library
if [ -e minicap/aosp/android-$rel/$abi/minicap.so ]; then
  adb -s $id push minicap/aosp/android-$rel/$abi/minicap.so $dir
else
  adb -s $id push minicap/aosp/android-$sdk/$abi/minicap.so $dir
fi

# Run!
adb -s $id shell LD_LIBRARY_PATH=$dir $dir/$bin $args "$@"

# Clean up
#adb shell rm -r $dir
