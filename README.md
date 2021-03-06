<p align="center">
<br>
<br>
<img src="https://sternentstehung.de/cleanup.png">
<br>
<br>
</p>

<h3 align="center">cleanup &ndash; Remove gone Git branches with ease.</h3>

<p align="center">
<a href="https://circleci.com/gh/dominikbraun/cleanup"><img src="https://circleci.com/gh/dominikbraun/cleanup.svg?style=shield"></a>
<a href="https://goreportcard.com/report/github.com/dominikbraun/cleanup"><img src="https://goreportcard.com/badge/github.com/dominikbraun/cleanup"></a>
<a href="https://www.codefactor.io/repository/github/dominikbraun/cleanup"><img src="https://www.codefactor.io/repository/github/dominikbraun/cleanup/badge" /></a>
<a href="https://github.com/dominikbraun/cleanup/releases"><img src="https://img.shields.io/github/v/release/dominikbraun/cleanup?sort=semver"></a>
<a href="https://github.com/dominikbraun/cleanup/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-Apache--2.0-brightgreen"></a>
<br>
<br>
<br>
</p>

---

:dizzy: cleanup is a CLI tool for keeping your Git repositories clean. It removes old branches in one or more repositories with a single command.

**Project status:** In active development. Can be used without warranty.

## <img src="https://sternentstehung.de/cleanup-dot.png"> Quick example

You merely need to provide the path to your repository. For example, change into your project directory and remove all gone branches like this:

````shell script
$ cleanup branches .
````

If you want to get a preview of branches that will be deleted, just perform a dry run.

````shell script
$ cleanup branches --dry-run .
````

There appear some branches that probably shouldn't be deleted? You can simply exclude them:

````shell script
$ cleanup branches --exclude="feature/3, feature/4" .
````

**Deleting branches that match a custom search term**

`--where` allows you to specify a string that has to be included in a branch's `git branch -vv` output. This option
overrides the default filter for gone branches, _meaning that also non-gone branches will be deleted_. 

For example, deleting _all_ feature branches is as simple as:

```shell script
$ cleanup branches --where="feature/" .
```

In case you want to delete _only gone feature branches_, extend the default filter for gone branches with `--and-where`.

```shell script
$ cleanup branches --and-where="feature/" .
```

**Cleaning multiple repositories at once**

Most developers work on multiple projects at once. Let's assume that these projects have a common parent directory.

<img src="https://sternentstehung.de/example-projects.png">

Deleting gone branches in all of these repositories is fairly simple:

````shell script
$ cleanup branches --has-multiple-repos projects
````

## <img src="https://sternentstehung.de/cleanup-dot.png"> What cleanup does

To find out which branches are gone on the remote, cleanup uses the output of `git branch -vv`.

````shell script
$ git branch -vv
* master	34a234a [origin/master] Merged some features
  feature/1	34a234a [origin/feature/1: gone] Implemented endpoints
  feature/2	3fc2e37 [origin/feature/2: gone] Added CLI flags
````

Based on that output, cleanup performs `git branch -d <branch>` on all gone branches. Maybe you've been doing the same
thing manually for each branch in each repository. However, the `--where` flag allows you to specify a custom string
that has to be included to delete the branch.

## <img src="https://sternentstehung.de/cleanup-dot.png"> Installation

Download the [latest release of cleanup](https://github.com/dominikbraun/cleanup/releases).

**Linux/macOS**

Copy the downloaded binary into a directory like `/usr/local/bin`. Make sure the directory is listed in `PATH`.

**Windows**

Create a directory like `C:\Program Files\cleanup` and copy the executable into it. [Add the directory](https://www.computerhope.com/issues/ch000549.htm) to `Path`.

## <img src="https://sternentstehung.de/cleanup-dot.png"> Remove branches periodically

Especially developers working in larger teams may want to clean up their repositories periodically.

**Linux/macOS**

Edit cron's job table.

````shell script
$ crontab -e
````

Adding a line like this will run cleanup every day at 11 PM for _all_ repositories. Make sure you provide correct paths to the cleanup binary and project directory.

````shell script
0 23 * * * /usr/local/bin/cleanup branches --has-multiple-repos /home/user/projects
````

**Windows**

Just create a scheduled task like this &ndash; and make sure to provide correct paths for the cleanup executable and project directory.

````shell script
> SCHTASKS /CREATE /SC DAILY /ST 23:00 /TN "Clean up repositories" ^
  /TR "C:\Path\To\cleanup.exe branches --has-multiple-repos C:\Path\To\Projects"
````

<p align="center">
<br>
<strong>クリーン</strong>
<br>
</p>