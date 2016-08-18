FROM flexconstructor/docker-centos-golang
MAINTAINER FlexConstructor <flexconstructor@gmail.com>

COPY go_app /go

RUN go-wrapper download github.com/shirou/gopsutil   \
                        github.com/sevlyar/go-daemon 
RUN go-wrapper install  github.com/shirou/gopsutil   \
                        github.com/sevlyar/go-daemon \
                        system
RUN mkdir /go/logs \
 && chmod -R 777 /go/logs


RUN echo  "[program:system_monitor]" >> /etc/supervisord.conf \
&& echo  "command = go run /go/src/sm_agent.go" >> /etc/supervisord.conf \
&& echo "startretries=1"  >> /etc/supervisord.conf
VOLUME /go/logs
#CMD ['go run /go/src/sm_agent.go']

ENTRYPOINT ["/usr/bin/supervisord","-n"]