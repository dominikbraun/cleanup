<p align="center">
<br>
<br>
<img src="https://sternentstehung.de/cleanup.png">
<br>
<br>
</p>

<h3 align="center">cleanup &ndash; Remove gone Git branches with ease.</h3>

<p align="center">
<img src="https://circleci.com/gh/dominikbraun/foodunit.svg?style=shield">
<img src="https://goreportcard.com/badge/github.com/dominikbraun/foodunit">
<img src="https://www.codefactor.io/repository/github/dominikbraun/dice/badge?s=0f13518b90c29be6bc3ec4ff537581a2e5c51c6a" />
<img src="https://img.shields.io/github/v/release/dominikbraun/foodunit?sort=semver">
<img src="https://img.shields.io/badge/license-Apache--2.0-brightgreen">
<br>
<br>
<br>
</p>

---

:dizzy: cleanup is a CLI tool for keeping your Git repositories clean. It removes old branches in one or more repositories with a single command.

**Project status:** In active development and stable to use.

## <img src="https://sternentstehung.de/cleanup-dot.png"> Usage

You merely need to provide the path to your repository. For example, change into your project directory and remove all gone branches:

````shell script
$ cleanup branches .
````

If you want to get a preview of branches that will be deleted, just perform a dry run.

````shell script
$ cleanup branches --dry-run .
````

There appear some branches that probably shouldn't be deleted? Exclude them:

````shell script
$ cleanup branches --exclude="feature/3, feature/4" .
````

#### Cleaning multiple repositories at once

Many developers are working on several projects. Let's assume that these projects have a common parent directory.

<img src="https://sternentstehung.de/example-projects.png">

Deleting gone branches in all of these repositories is fairly simple:

````shell script
$ cleanup branches --has-multiple-repos projects
````

## <img src="https://sternentstehung.de/cleanup-dot.png"> Installation

Download the [latest release of cleanup](https://github.com/dominikbraun/cleanup/releases) and copy the binary into a directory like `/usr/local/bin`. Make sure the directory is listed in `PATH`.

## <img src="https://sternentstehung.de/cleanup-dot.png"> Remove branches periodically

##### Linux

...

##### macOS

...

##### Windows

...