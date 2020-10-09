#!/usr/bin/env bash

RUN_FILE_NAME=reskd

run(){
	echo running... ./${RUN_FILE_NAME} $2
	./${RUN_FILE_NAME} $2
}

start(){
	echo 'starting $2'
	nohub ./${RUN_FILE_NAME} $2 > ../logs/resk.log 2>&1 &
	echo 'started ${RUN_FILE_NAME} $2'
}

stop(){
	pids=`ps -ef|grep ${RUN_FILE_NAME}|grep -v 'grep'|awk '{print $2}'`
	echo '$pids'
	for pid in $pids
	do
		echo kill $pid
		kill -15 $pid
	done
}

restart(){
	stop && start $@
}

rertun(){
	stop && run $@
}

#./run.sh run 参数
action='$1'
if ['$action' == '']
then
	action='run'
fi

case '${action}' in
	start)
		start '$@';;
	stop)
		stop '$@';;
	restart)
		restart '$@';;
	run)
		run '$@';;
	rerun)
		retun '$@';;
	*)
		echo 'Usage: $0 {start|stop|restart|run|rerun} {dev|test|prod|...}';
		echo '		eg: ./${RUN_FILE_NAME} run dev';
		exit 1;
esac

exit0