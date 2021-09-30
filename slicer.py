import xml.etree.ElementTree as ET
import argparse
from os import path
import json


xmlns = "{http://schemas.android.com/apk/res/android}"
analysis = dict()


def is_accessible(child):
    """check if the component is accessible or not.

    Args:
        child : children element of the component

    Returns:
        activities(dict): dictionary format containing
                          information about the component
    """
    activities = {}
    filters = {}
    if xmlns+"exported" in child.attrib:
        # It's exported so just get the intents and play with it
        pass
    
    if child.find("intent-filter"):
        name = child.attrib[xmlns+"name"]
        
        for element in child.findall("intent-filter"):
            for elem in element:
                # Here is a possibility of error like out of index
                filters[elem.tag] = list(elem.attrib.values())[0]
        activities[child.tag] = name
        activities['intent-filters'] = filters

    return activities


def process_manifest(tree):
    """process the Android manifest and find out all
     the exported or accessible components

    Args:
        tree: xml.ElementTree Object
    """

    for child in tree.iter():
        if child.tag == "manifest":
            analysis['version'] = child.attrib[xmlns+'versionName']

            if xmlns+'allowBackup' in child.attrib:
                analysis['allowBackup'] = child.attrib[xmlns+'allowBackup']
            analysis['package'] = child.attrib['package']

        elif child.tag == "application":
            if xmlns+"debuggable" in child.attrib:
                analysis['debuggable'] = child.attrib[xmlns+"debuggable"]
        
        elif child.tag in ["activity", "service", "receiver", "provider"]:
            activity = is_accessible(child)
            if activity:
                print(activity)


def main(directory: str, config_file: str):
    """Drive the whole program

    Args:
        directory (str): directory which we have to process
        config_file (str): path to the config file
    """
    if path.isdir(directory) and path.exists(config_file):
        with open(config_file, 'r') as f:
            config = json.load(f)

    manifest_path = path.join(directory, config['paths']['manifest'])
    if path.isfile(manifest_path):
        tree = ET.parse(manifest_path)
        process_manifest(tree)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="pubg ranking for all servers")
    parser.add_argument("-d", "--dir", help="path to the jadx directory")
    parser.add_argument("-c", "--config", help="path to the config file in json format")
    args = parser.parse_args()

    if args.dir and args.config:
        main(args.dir, args.config)
