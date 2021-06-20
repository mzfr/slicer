#![allow(non_snake_case, unused_must_use)]
use clap::{App, Arg, ArgMatches};
use quickxml_to_serde::{xml_string_to_json, Config};
use std::path::Path;

static MANIFEST: &str = "resources/AndroidManifest.xml";
// static XMLNS: &str = "http://schemas.android.com/apk/res/android";

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

fn general_package_info(jsondata: serde_json::Value) {}
// Given a directory name it will read the AndroidManifest.xml
// file and return the XML document loaded in the memory
fn read_xml_file(directory: String) {
    if Path::new(&directory).exists() {
        // Join paths and then store the string version into another variable
        let AndroidManifestPath = Path::new(&directory).join(MANIFEST); // This gives PathBuf
        let xmlpath = AndroidManifestPath.display().to_string();

        let text = std::fs::read_to_string(xmlpath).unwrap();
        let conf = Config::new_with_defaults();
        let jsondata = xml_string_to_json(text.to_owned(), &conf);

        //for (k, v) in json.iter().enumerate() {
        //println!("{:?} ---- {:?}\n", k, v);
        //}
        general_package_info(jsondata);
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
