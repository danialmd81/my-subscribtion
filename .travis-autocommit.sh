#!/bin/bash
set -e

# Configure Git
git config --global user.name "Travis CI"
git config --global user.email "travis@users.noreply.github.com"

# Add & Commit changes
git add .
git commit -m "Automated update from Travis CI [skip ci]" || exit 0 # Skip if no changes

# Push using a GitHub Personal Access Token (PAT)
git remote set-url origin https://${GH_TOKEN}@github.com/${TRAVIS_REPO_SLUG}.git
git push origin HEAD:${TRAVIS_BRANCH}