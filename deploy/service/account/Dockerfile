FROM centos:7

ADD bin/account /

RUN chmod 777 /account

ENV PARAMS=""

ENTRYPOINT ["sh","-c","/account $PARAMS"]