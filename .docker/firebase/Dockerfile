FROM node:18.15.0-alpine

WORKDIR /firebase

COPY firebase/firebase.json /firebase/

COPY firebase/.firebaserc /firebase/

RUN npm install -g firebase-tools

EXPOSE 9099

EXPOSE 4000

CMD ["firebase", "emulators:start", "--project=emulator"]
