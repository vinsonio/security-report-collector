# Contributing to Security Report Collector

First off, thank you for considering contributing to Security Report Collector! It's people like you that make open source such a great community.

## Where do I go from here?

If you've noticed a bug or have a feature request, [make one](https://github.com/vinsonio/security-report-collector/issues/new)! It's generally best if you get confirmation of your bug or approval for your feature request this way before starting to code.

### Fork & create a branch

If this is something you think you can fix, then fork Security Report Collector and create a branch with a descriptive name.

A good branch name would be (where issue #325 is the ticket you're working on):

```sh
git checkout -b add-mongodb-support
```

### Get the test suite running

Make sure you're able to run the test suite locally. You'll need Docker and Docker Compose installed.

```sh
go test ./...
```

### Implement your fix or feature

At this point, you're ready to make your changes! Feel free to ask for help; everyone is a beginner at first :smile_cat:

### Make a Pull Request

At this point, you should switch back to your master branch and make sure it's up to date with Security Report Collector's master branch:

```sh
git remote add upstream git@github.com:vinsonio/security-report-collector.git
git checkout main
git pull upstream main
```

Then update your feature branch from your local copy of master, and push it!

```sh
git checkout add-mongodb-support
git rebase main
git push --set-upstream origin add-mongodb-support
```

Finally, go to GitHub and make a Pull Request.

### Keeping your Pull Request updated

If a maintainer asks you to "rebase" your PR, they're saying that a lot of code has changed, and that you need to update your branch so it's easier to merge.

To learn more about rebasing, check out this guide: https://www.atlassian.com/git/tutorials/rewriting-history/git-rebase

## Code of Conduct

We have a Code of Conduct, please follow it in all your interactions with the project.