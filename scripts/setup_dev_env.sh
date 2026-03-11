export APP_ENV="DEV"
export SYSTEM_EXTERNAL_URL="https://localhost"
export JWT_SECRET="devsecret"
export INTERNAL_AUTH_SECRET="devsecret"
export FILESTORAGE_ACCESS_KEY="devsecret"
export FILESTORAGE_SECRET_KEY="devsecret"
export DISCORD_BOT_ADMIN_USER_ID="1081107248417361970"
export POSTGRES_PASSWORD="devsecret"
export MINECRAFT_RCON_PASSWORD="devsecret"
export ALERTS_DISCORD_WEBHOOK=""
export GF_SECURITY_ADMIN_USER="devsecret"
export GF_SECURITY_ADMIN_PASSWORD="devsecret"
export CONTAINER_REGISTRY="ghcr.io"
export GIT_FULL_REPO_NAME=`git remote get-url origin | awk -F: '{sub(/\.git$/, "", $2); print $2}'`
export DOCKER_IMAGE_TAG="main"
export CERTBOT_EMAIL="contact@varso.org"
export CERTBOT_DOMAINS="varso.org"
export USE_TEST_CERT="true"
export NPM_REGISTRY_URL="ghcr.io"
export AUTH_ADMIN_EMAIL="contact@varso.org"
export AUTH_ADMIN_PASSWORD="devsecret"
export PROD_EC2_PROJECT_PATH="."
export PLAUSIBLE_SECRET_KEY="devsecretdevsecretdevsecretdevsecretdevsecretdevsecretdevsecretdevsecret"

echo "NPM Token (optional, used for building www production image): "
read npm_token
export NPM_TOKEN=$npm_token 

./scripts/apply_env.sh .env
./scripts/apply_env.sh ./src/services/www/src/config/config.ts
