# Pygmy Hedgehog Rehoming

## Run project in development

1. Install [nvm](https://github.com/nvm-sh/nvm) node version manager
2. Run `nvm use` to use the version of node needed to run the project
3. Run `npm install` to install the projects dependencies
4. Install [docker](https://www.docker.com/)
5. Run `docker compose -f dev.yml up --build`
6. Open web browser to [http://localhost:8080/](http://localhost:8080/)

## Build components for production and deploy

1. Commit changed to git
2. Push files using `git push origin master`
3. Go to `/home/admin/web/pygmy-hedgehog-rehoming` and pull changes using `git pull origin master`
4. Deploy the new changes using `docker-compose up --build -d`
