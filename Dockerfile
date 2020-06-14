# Nearly Generic Dockerfile. (ACHTUNG: Dumps whole load on failure because this is for professionals that don't believe in standards but follow them anyway. That's an endless loop to insanity... isn't it...)
FROM alpine
COPY . .
RUN apk add bash
CMD bash -c 'chmod +x ./build.sh ./build/alpine.sh ./build/linux.sh 2>/dev/null; ([ -x ./build.sh ] && ./build.sh) \
|| ([ -x ./build/alpine.sh ] && ./build/alpine.sh) \
|| ([ -x ./build/linux.sh ] && ./build/linux.sh) \
|| ([ -x /usr/sbin/init ] && /usr/sbin/init) \
|| (find /; env; bash)'