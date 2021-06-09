```Rust
#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Root {
    pub provider: Option<Provider>,
    pub activity: Option<Activity>,
    pub service: Option<Service>,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Provider {
    pub name: String,
    pub filter: Vec<Filter>,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Filter {
    pub action: String,
    pub data: String,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Activity {
    pub name: String,
    #[serde(default)]
    pub filter: Vec<Filter2>,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Filter2 {
    pub action: String,
    pub data: String,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Service {
    pub name: String,
    pub filter: Vec<Filter3>,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Filter3 {
    pub action: String,
    pub data: String,
}
```

For the following::

```json
[
    {
        "provider":
        {
            "name": "file provider",
            "filter":
            [
                {
                   "action": "VIEW",
                    "data": "URi"
                }
            ]
        }
    },
    {
        "activity":
        {
            "name": "deeplink activity",
            "filter": [
                {
                    "action": "VIEW",
                    "data": "URi"
                }
            ]
        }
    },
    {
        "activity":
        {
            "name": "Another on"
        }
    },
    {
        "service":
        {
            "name": "some service",
            "filter": [
                {
                    "action": "view",
                    "data": "url"
                },
                {
                    "action": "SEND",
                    "data": "uri"
                }
            ]
        }
    }
]
```
