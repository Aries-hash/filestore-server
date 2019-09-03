FROM centos:7

ADD bin/apigw /

RUN chmod 777 /apigw

ENV PARAMS=""

ENTRYPOINT ["sh","-c","/apigw $PARAMS"]