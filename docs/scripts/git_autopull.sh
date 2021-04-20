#!/bin/bash

# Arthur 84500316@qq.com
# 背景：平时收集的工具比较多，比赛或hvv能够批量更新这些小工具。
# 原理：遍历目录发现.git强制覆盖更新，更新后返回上一层目录继续查找

ignoreFiles=".DS_Store . .."

list() {
	find . -maxdepth 1 
}

ignore() {
	fi=$1
	for f in $ignoreFiles; do
		if [ `basename "$fi"` = "$f" ]; then
			echo "yes"
			return
		fi	
	done
	echo "no"
}

isInTargetDir() {
	files=`list`
	for f in $files; do
		if [ `basename "$f"` = ".git" ]; then
			echo "yes"
			return
		fi
	done
	echo "no"
	return
}

check() {
	ok=`isInTargetDir`
	if [ "$ok" = "yes" ]; then
		echo "update" "$PWD"
		 git fetch --all &&  git reset --hard origin/master && git pull
	else
		files=`list`
		for f in $files; do
			if [ ! -d "$f" ]; then
			       continue
			fi       
			ok=`ignore "$f"`
			if [ "$ok" = "yes" ]; then 
			       	continue
			fi
			cd $f
			check
			cd ../ 
		done
	fi
}

check
