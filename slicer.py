import xml.etree.ElementTree as ET
import argparse
from os import path
import requests
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
    analysis["Google Keys"] = analysis["Keys/Tokens"] = []
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
            component = is_accessible(child)
            if component:
                if child.tag not in analysis:
                    analysis[child.tag] = []
                analysis[child.tag].append(component)


def find_keys(strings_path: str, config):
    """Find the api keys in strings.xml and AndroidManifest.xml

    Args:
        strings_path (str): path to the strings.xml file
        config (json): JSON object of the config file
    """

    values = ET.parse(strings_path)
    for child in values.iter():
        if "name" in child.attrib:
            attrib_name = child.attrib['name']
            if attrib_name == "firebase_database_url":
                url = child.text + "/.json"
                r = requests.get(url)
                if r.status_code != 401 and "disabled" not in r.content:
                    analysis['Firebase'] = (url, r.status_code)

            elif attrib_name == "google_api_key" or attrib_name == "google_map_keys":
                key = child.text
                for _, v in config['URLs'].items():
                    url = v + key
                    r = requests.get(url)
                    if r.status_code != 403 and b"API project is not authorized" not in r.content:
                        analysis["Google Keys"].append({url: r.status_code})
            else:
                for i in ["api", "keys", "token"]:
                    if i in attrib_name.lower() and attrib_name.lower() not in ["abc_capital_off", "abc_capital_on", "currentApiLevel"]:
                        analysis['Keys/Tokens'].append(
                            {attrib_name: child.text})


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
    strings_path = path.join(directory, config['paths']['strings'])

    if path.isfile(manifest_path):
        tree = ET.parse(manifest_path)
        process_manifest(tree)
        find_keys(strings_path, config)

    with open(analysis['package']+'.json', "w") as f:
        json.dump(analysis, f)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Slicer")
    parser.add_argument("-d", "--dir", help="path to the jadx directory")
    parser.add_argument(
        "-c", "--config", help="path to the config file in json format")
    args = parser.parse_args()

    if args.dir and args.config:
        main(args.dir, args.config)
