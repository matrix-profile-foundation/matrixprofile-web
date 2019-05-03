FROM node:8-alpine
RUN ["npm", "install", "-g", "serve"]
USER node
WORKDIR /home/node/app
ENV NODE_ENV production
ADD . /home/node/app
EXPOSE 8080
CMD ["serve", "-s", "dist", "-l", "8080"]
