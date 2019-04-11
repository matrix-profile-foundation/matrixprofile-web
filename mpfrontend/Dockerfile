FROM node:8-alpine
USER node
WORKDIR /home/node/app
ENV NODE_ENV production
ADD . /home/node/app
EXPOSE 8080
CMD ["npm", "run", "serve"]
