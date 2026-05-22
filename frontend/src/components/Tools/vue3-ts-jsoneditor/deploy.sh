#!/usr/bin/env sh

# final commit
git add -A
git commit -m 'deploy'

# abort on errors
set -e

# update package version
npm version "$1"

# build
npm run build
npm run build:gh

# publish
npm publish

# navigate into the build output directory
cd gh-dist

# if you are deploying to a custom domain
# echo 'www.example.com' > CNAME

git init
git checkout -b main
git add -A
git commit -m 'deploy'

# if you are deploying to https://<USERNAME>.github.io/<REPO>
git push -f git@github.com:bestkolobok/vue3-jsoneditor.git main:gh-pages

rm -rf .git

cd -