[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![platform](https://img.shields.io/badge/platform-osx%2Flinux%2Fwindows-green.svg)](https://github.com/mzfr/slicer)

# Slicer

A tool to automate the recon process on an APK file. 

Slicer accepts a path to an extracted APK file and then returns all the activities, receivers, and services which are exported and have `null` permissions and can be externally provoked.

__Note__: The APK has to be extracted via `jadx` or `apktool`.

# Table of Content

- [Slicer](#slicer)
- [Table of Content](#table-of-content)
- [Summary](#summary)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Usage Example](#usage-example)
- [Acknowledgements and Credits](#acknowledgements-and-credits)
- [Contribution](#contribution)
- [Support](#support)

# Summary

__Why?__

I started bug bounty like 3 weeks ago(in June 2020) and I have been trying my best on android apps. But I noticed one thing that in all the apps there were certain things which I have to do before diving in deep. So I just thought it would be nice to automate that process with a simple tool. 

__Why not drozer?__

Well, drozer is a different beast. Even though it does finds out all the accessible components but I was tired of running those commands again and again.

__Why not automate using drozer?__

I actually wrote a bash script for running certain drozer commands so I won't have to run them manually but there was still some boring stuff that had to be done. Like Checking the `strings.xml` for various API keys, testing if firebase DB was publically accessible or if those google API keys have setup any cap or anything on their usage and lot of other stuff.

__Why not search all the files?__

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

* Clone this repository

```
git clone https://github.com/mzfr/slicer
```
* `cd slicer`
* Now you can run it: `python3 slicer.py -h`

# Usage

It's very simple to use. Following options are available:

```
Extract information from Manifest and strings of an APK

Usage:
        slicer [OPTION] [Extracted APK directory]

Options:

  -d, --dir             path to jadx output directory
  -o, --output          Name of the output file(not implemented)
```

I have not implemented the `output` flag yet because I think if you can redirect slicer output to a yaml file it will a proper format.

# Usage Example

* Extract information from the APK and display it on the screen.

```bash
python3 slicer.py -d path/to/extact/apk -c config.json
```

# Acknowledgements and Credits

The extractor module used to extract URLs and paths is taken from [apkurlgrep](https://github.com/ndelphit) by @ndelphit

# Contribution

All the features implemented in this are things that I've learned in past few weeks, so if you think that there are various other things which should be checked in an APK then please open an issue for that feature and I'd be happy to implement that :)

# Support

If you'd like you can buy me some coffee:

<a href="https://www.buymeacoffee.com/mzfr" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" style="height: 51px !important;width: 217px !important;" ></a>
