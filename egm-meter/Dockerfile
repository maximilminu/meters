FROM scratch

EXPOSE 80
WORKDIR /

ENV ZONEINFO=/zoneinfo.zip
COPY zoneinfo.zip /

COPY egm-meter /

CMD ["/egm-meter"]
