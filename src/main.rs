#![allow(non_snake_case, unused_must_use)]
use quickxml_to_serde::{xml_string_to_json, Config};


fn main() {     
    let text = std::fs::read_to_string("D:\\dev\\slicer\\src\\Manifest.xml").unwrap();
    let conf = Config::new_with_defaults();
    let json = xml_string_to_json(text.to_owned(), &conf);
    //println!("{:?}", json.unwrap()["manifest"]);

    for (k,v) in json.unwrap()["manifest"].it {
        println!("{:?} ---- {:?}\n", k, v);
    }
}
