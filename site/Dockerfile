FROM node:slim

ENV ENV production

# Setting up the work directory
WORKDIR /app

# Copying all the files in our project
COPY . .

# Installing dependencies
RUN npm install

# Starting our application
CMD [ "node", "./bin/www" ]

# Exposing server port
EXPOSE 5000