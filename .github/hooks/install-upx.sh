#!/usr/bin/env bash

[ -f /usr/bin/upx ] && exit

function Init_Arch() {
  ARCH=$(uname -m)
  case $ARCH in
  armv5*) ARCH="armv5" ;;
  armv6*) ARCH="armv6" ;;
  armv7*) ARCH="arm" ;;
  aarch64) ARCH="arm64" ;;
  x86) ARCH="386" ;;
  x86_64) ARCH="amd64" ;;
  i686) ARCH="386" ;;
  i386) ARCH="386" ;;
  esac
}

Init_Arch

upx_version=$(wget -qO- -t5 -T10 "https://api.github.com/repos/upx/upx/releases/latest" |
  grep "tag_name" | head -n 1 | awk -F ":" '{print $2}' | sed 's/\"//g;s/,//g;s/ //g')

cd /tmp/
wget https://github.com/upx/upx/releases/download/${upx_version}/upx-${upx_version/v/}-${ARCH}_linux.tar.xz -O \
  upx-${upx_version/v/}-${ARCH}_linux.tar.xz || echo '下载 upx 失败!'
tar -xvf upx-${upx_version/v/}-${ARCH}_linux.tar.xz
cd upx-${upx_version/v/}-${ARCH}_linux
\mv upx /usr/bin/

# 输出安装完成信息
upx --version
echo 'Upx 安装完成！'
