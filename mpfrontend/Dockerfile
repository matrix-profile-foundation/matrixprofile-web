FROM node:8-alpine
RUN ["npm", "install", "-g", "serve"]
WORKDIR /home/node/app
USER node
EXPOSE 8080
CMD ["serve", "-s", "dist", "-l", "8080"]
