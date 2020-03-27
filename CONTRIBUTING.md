# Contributing to rundeck-zabbix

First off, thank you for considering contributing to rundeck-zabbix :tada:. Following these guidelines help to communicate that you respect the time of the developers managing and developing this open source project. In return, they should reciprocate that respect in addressing your issue, assessing changes, and helping you finalise your pull requests.

### Where do I go from here?

If you've noticed a bug or have a question, [search the issue tracker](https://github.com/kosyfrances/rundeck-zabbix/issues?q=something) to see if
someone else has already created a ticket. If not, go ahead and
[make one](https://github.com/kosyfrances/rundeck-zabbix/issues/new)!

### Fork & create a branch

If this is something you think you can fix, then [fork rundeck-zabbix](https://help.github.com/en/articles/fork-a-repo) and
create a branch with a descriptive name.

A good branch name would be (where issue #225 is the ticket you're working on):

```sh
git checkout -b 225-add-response
```

### Set up dev environment
Refer to the [dev setup guide](/dev/README.md).


### Implement your fix or feature

At this point, you're ready to make your changes! Feel free to ask for help;
everyone is a beginner at first.

This project uses [go modules](https://blog.golang.org/using-go-modules) to manage dependencies.

Running `make test` runs tests.

Also, remember to write tests :smile_cat: We care
about quality, so your PR won't be merged until all tests pass.

### Keeping your Pull Request updated

If a maintainer asks you to "rebase" your PR, they're saying that a lot of code
has changed, and that you need to update your branch so it's easier to merge.

### Merging a PR (maintainers only)

A PR can only be merged into master by a maintainer if:

* It is passing CI.
* It has no requested changes.
* It is up to date with current master.

Any maintainer is allowed to merge a PR if all of these conditions are
met.

### Shipping a release (maintainers only)
Releases are done on the [Github UI](https://github.com/kosyfrances/rundeck-zabbix/releases) using [SemVer](https://semver.org/) for versioning.
