FROM node:19.8.1-alpine

WORKDIR /firebase

RUN npm install -g firebase-tools

EXPOSE 9099

EXPOSE 4000

CMD ["firebase", "emulators:start", "--project=emulator"]
