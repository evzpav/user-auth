# Dokku Deployment

- Add .gitlab-ci.yml:
```yaml
stages:
  - deploy
  
variables:
  APP_NAME: user-auth
  
deploy:
  stage: deploy
  image: ilyasemenov/gitlab-ci-git-push
  environment:
    name: production
  only:
    - master
  script:
    - git-push ssh://dokku@$VM_IP:22/$APP_NAME
```

- Create project on Gitlab
- Push project to master
- Set env vars VM_IP and SSH_PRIVATE_KEY on Gitlab

## On VM machine:
```bash
# Set MySQL:
sudo dokku plugin:install https://github.com/dokku/dokku-mysql.git mysql

export PROJECT_NAME=user-auth

export MYSQL_IMAGE="mysql"
export MYSQL_IMAGE_VERSION="8.0"

dokku mysql:create $PROJECT_NAME

dokku mysql:link $PROJECT_NAME $PROJECT_NAME # It will set DATABASE_URL env var in the project


dokku mysql:connect $PROJECT_NAME 
#Run scripts form docker/mysql/init.sql

```