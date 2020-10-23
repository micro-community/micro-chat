FROM alpine
ADD chat /chat
ENTRYPOINT [ "/chat" ]
