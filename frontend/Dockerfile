FROM node:lts-alpine
WORKDIR /usr/src/app
RUN npm install -g pnpm
COPY ["package.json", "package-lock.json*", "npm-shrinkwrap.json*", "./"]
RUN pnpm install && pnpm add vite
COPY . .
# ENV NODE_ENV=production
# RUN pnpm build
EXPOSE 80
# RUN chown -R node /usr/src/app
# USER node
CMD ["pnpm", "start"]
