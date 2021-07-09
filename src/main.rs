#![allow(non_snake_case, unused_must_use)]
use clap::{App, Arg, ArgMatches, AppSettings};
use quickxml_to_serde::{xml_string_to_json, Config};
use std::collections::HashMap;
use std::process;
use serde_json::Value;
use std::path::Path;

static MANIFEST: &str = "resources/AndroidManifest.xml";
// static XMLNS: &str = "http://schemas.android.com/apk/res/android";

fn parse_args() -> ArgMatches {
    App::new("slicer")
        .setting(AppSettings::ArgRequiredElseHelp)
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

fn general_package_info(jsondata: serde_json::Value) -> HashMap<&'static str, &'static Value>{
    let mut package_info = HashMap::new();
    for (k,v) in jsondata["manifest"].as_object().unwrap() {
        if k == "@android:allowBackup" {
            package_info.insert("backup", v);
        } else if k == "@package" {
            package_info.insert("name", v);
        } else if k == "@android:versionName" {
            package_info.insert("version", v);
        } else if k == "@android:debuggable" {
            package_info.insert("debug", v);
        }
    }

    return package_info;

}
// Given a directory name it will read the AndroidManifest.xml
// file and return the XML document loaded in the memory
fn read_xml_file(directory: String) -> serde_json::Value {
    let jsondata;
    if Path::new(&directory).exists() {
        // Join paths and then store the string version into another variable
        let AndroidManifestPath = Path::new(&directory).join(MANIFEST); // This gives PathBuf
        let xmlpath = AndroidManifestPath.display().to_string();

        let text = std::fs::read_to_string(xmlpath).unwrap();
        let conf = Config::new_with_defaults();
        jsondata = xml_string_to_json(text.to_owned(), &conf);

    } else {
        println!("Failed to Read the file");
        process::exit(1);
    } 
    
    return jsondata.unwrap();
}

fn main() {
    let args = parse_args();

    let directory = args.value_of("dir").unwrap();
    let jsondata = read_xml_file(directory.to_string());

    let package_info = general_package_info(jsondata);
}
