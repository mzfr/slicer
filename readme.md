# Slicer

A tool to automate the recon process on an APK file. 

Slicer accepts a path to an extracted APK file and then returns all the activities, receivers, and services which are exported and have `null` permissions and can be externally provoked.

__Note__: The APK has to be extracted via `jadx` or `apktool`.

# Languages

So I initially wrote slicer in `golang` but then I just wanted to learn `rust` so I decided to port the tool in golang. If you'd like to see the golang code you can checkout the `go` branch. The master branch contains only the rust files.

P.S -> As of 26/05/2021, rust version of slicer is not supposed to be used because it's not fully functional and I'm still working on it.

# Summary

### Why?

I started bug bounty like 3 weeks ago(in June 2020) and I have been trying my best on android apps. But I noticed one thing that in all the apps there were certain things which I have to do before diving in deep. So I just thought it would be nice to automate that process with a simple tool. 
### Why not drozer?

Well, drozer is a different beast. Even though it does finds out all the accessible components but I was tired of running those commands again and again.

### Why not automate using drozer?

I actually wrote a bash script for running certain drozer commands so I won't have to run them manually but there was still some boring stuff that had to be done. Like Checking the `strings.xml` for various API keys, testing if firebase DB was publically accessible or if those google API keys have setup any cap or anything on their usage and lot of other stuff.

### Why not search all the files?

I think that a tool like grep or ripgrep would be much faster to search through all the files. So if there is something specific that you want to search it would be better to use those tools. But if you think that there is something which should be checked in all the android files then feel free to open an issue.

# Features

* Check if the APK has set the `android:allowbackup` to `true`
* Check if the APK has set the `android:debuggable` to `true`.
* Return all the activities, services and broadcast receivers which are exported and have null permission set. This is decided on the basis of two things:
    - `android:exporte=true` is present in any of the component and have no permission set.
    -  If exported is not mention then slicer check if any `Intent-filters` are defined for that component, if yes that means that component is exported by default(This is the rule given in android documentation.)

* Check the Firebase URL of the APK by testing it for `.json` trick.
    - If the firebase URL is `myapp.firebaseio.com` then it will check if `https://myapp.firebaseio.com/.json` returns something or gives permission denied.
    - If this thing is open then that can be reported as high severity.

* Check if the google API keys are publically accessible or not. 
    - This can be reported on some bounty programs but have a low severity.
    - But most of the time reporting this kind of thing will bring out the pain of `Duplicate`.
    - Also sometimes the company can just close it as `not applicable` and will claim that the KEY has a `usage cap` - r/suspiciouslyspecific :wink: 

* Return other API keys that are present in `strings.xml` and in `AndroidManifest.xml`
* List all the file names present in `/res/raw` and `res/xml` directory.
* Extracts all the URLs and paths.
    - These can be used with tool like dirsearch or ffuf.


# Installation

## For Rust

If you'd like to build this locally, then make sure you have [`rust` installed](https://www.rust-lang.org/tools/install). After that do the following

* `git clone https://github.com/mzfr/slicer`
* `cd slicer`
* `cargo build`

And then in the `target/debug` there should be a binary named `Slicer` which you can run.

Once its ported completely in rust then I'll just release the binaries as well.

## For golang

__I am not sure if go get will work now since I've moved stuff to another branch.__

You can download the binary from the [release](https://github.com/mzfr/slicer/releases) page. Also if you want you can clone this repository and build the binary yourself.

If you have `go` compiler installed then you can use `go get github.com/mzfr/slicer`.

__NOTE__: Slicer uses `config.yml` file. So either have a file named `config.yml` in your current working directory or make a directory
named `.slicer` in your `$HOME` and then place the `config.yml` file there.

## Arch Linux

`slicer` can be installed from available [AUR packages](https://aur.archlinux.org/packages/?O=0&SeB=nd&K=A+tool+to+automate+the+boring&outdated=&SB=n&SO=a&PP=50&do_Search=Go) using an [AUR helper](https://wiki.archlinux.org/index.php/AUR_helpers). For example,

```
yay -S slicer
```

If you prefer, you can clone the [AUR packages](https://aur.archlinux.org/packages/?O=0&SeB=nd&K=A+tool+to+automate+the+boring&outdated=&SB=n&SO=a&PP=50&do_Search=Go) and then compile them with [makepkg](https://wiki.archlinux.org/index.php/Makepkg). For example,

```
git clone https://aur.archlinux.org/slicer.git && cd slicer && makepkg -si
```

# Usage

It's very simple to use. Following options are available:

```
Extract information from Manifest and strings of an APK

Usage:
        slicer [OPTION] [Extracted APK directory]

Options:

  -d, --dir             path to jadx output directory
 -nb, --no-banner       Don't Show Banner
```

# Usage Example

* Extract information from the APK and display it on the screen.

```
slicer -d path/to/extact/apk
```

* Extract information and store in a yaml file:

```
slicer -d path/to/extracted/apk -nb=false > name.yaml
```
__If you plan to use if for Bug bounty or anything similar it's better to store in some file__

# Contribution

All the features implemented in this are things that I've learned in past few weeks, so if you think that there are various other things which should be checked in an APK then please open an issue for that feature and I'd be happy to implement that :)

# Support

If you'd like you can buy me some coffee:

<a href="https://www.buymeacoffee.com/mzfr" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" style="height: 51px !important;width: 217px !important;" ></a>
