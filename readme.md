# Slicer

Okay so I am not going to follow the old way of implementing this. Instead I would make it a very simple tool and then later might add new functionalities to it.

The things that I have in mind right now are:

1) Read strings.xml and find all the api keys and firebase URL
    - For searching just pickup all the strings that have `api`/`token` in them.
        + Later we'll do a regex check if the values have format of an API or something I guess.
    - For Firebase perform the `/.json` check directly.
    - Maybe for google API key add those checks to see if the key is accessible or not.

2) It will list all the files of the directly `/res/raw/` and `/res/xml`
3) It will parse the androidmanifest.xml file and will show all the exported activites, broadcast, service and content.
    - Only show the one that have exported=True
        + It would be nice if we can have it like: if no `exported` mention then check if any `intent-filter` is defined. If yes then show that else don't.
    - Show all the intent filter of all the exported surfaces.
    - Also show the type of data URL they accept
    - check if the backup is true or false.
    - also check if the grantUriPermission is setup on any content provider link


The 3rd point looks like a work of drozer but if we can do this then I think we won't have to depend on drozer at all. Also having all this automated would really save loads of time
