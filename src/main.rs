#![allow(non_snake_case, unused_must_use)]
use clap::{App, Arg, ArgMatches};
use std::collections::HashMap;
use std::path::Path;

static MANIFEST: &str = "resources/AndroidManifest.xml";
static XMLNS: &str = "{http://schemas.android.com/apk/res/android}";

fn parse_args() -> ArgMatches {
    App::new("slicer")
        .version("2.0")
        .author("mzfr")
        .about("Automate boring process of APK recon")
        .args(&[Arg::new("dir")
            .about("directory path for the extracted APK")
            .short('d')
            .long("dir")
            .takes_value(true)])
        .get_matches()
}

//Get package name, version and other general things
fn general_package_info(node: roxmltree::Node) {
    let mut package_info = HashMap::new();

    for attrib in node.attributes() {
        //Improvement: In case these attributes doesn't exists
        // it will not show those values so if we want that then
        // instead of match use if/else
        match attrib.name() {
            "package" => package_info.insert("Package", attrib.value()),
            "versionName" => package_info.insert("Version", attrib.value()),
            "allowBackup" => package_info.insert("Allow Backup", attrib.value()),
            _ => None,
        };
    }
    //TODO: Return the hashmap
}

// Just check for {"activity", "receiver", "service", "provider"}
// Check if the component is exported, if yes then cool
// if no then check if it has intent filters or not
// if yes then parse over those filters as well
// if no then not exported
fn exported_components(node: roxmltree::Node) {
    let exported = format!("{}{}", XMLNS, "exported");
    let name = format!("{}{}", XMLNS, "name");
    if node.has_attribute("android:exported") {
        println!("It's working");
        //println!("{:?}", node.attribute(name.as_str()));
    } else {
        println!("Not working");
    }
}

// Given a directory name it will read the AndroidManifest.xml
// file and return the XML document loaded in the memory
fn read_xml_file(directory: String) {
    if Path::new(&directory).exists() {
        // Join paths and then store the string version into another variable
        let AndroidManifestPath = Path::new(&directory).join(MANIFEST); // This gives PathBuf
        let xmlpath = AndroidManifestPath.display().to_string();

        // TODO: Check wether the path exist or not
        // right now it will panic if the file is not found
        // maybe follow the similar pattern as below
        let text = std::fs::read_to_string(&xmlpath).unwrap();
        let doc = match roxmltree::Document::parse(&text) {
            Ok(doc) => doc,
            Err(e) => {
                println!("Error: {}.", e);
                return;
            }
        };
        for node in doc.descendants() {
            match node.tag_name().name() {
                "manifest" => general_package_info(node),
                "activity" => exported_components(node),
                _ => (), // In case the node is something different
            }
        }
    }
}

fn main() {
    let args = parse_args();

    if args.is_present("dir") {
        let directory = args.value_of("dir").unwrap();
        //let doc = read_xml_file(directory.to_string());
        read_xml_file(directory.to_string());
    }
}
