FROM centos:7

ADD bin/dbproxy /

RUN chmod 777 /dbproxy

ENV PARAMS=""

ENTRYPOINT ["sh","-c","/dbproxy $PARAMS"]