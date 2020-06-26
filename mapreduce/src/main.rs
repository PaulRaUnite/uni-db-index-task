use std::error::Error;
use std::io;

use chrono::{Date, DateTime, Utc};
use rayon::prelude::*;
use serde::Deserialize;
use std::collections::{BTreeMap, HashMap};
use std::hash::Hash;

#[derive(Debug, Deserialize)]
struct Record {
    #[serde(rename = "InvoiceNo")]
    invoice_no: u64,
    #[serde(rename = "StockCode")]
    stock_code: String,
    #[serde(rename = "Description")]
    description: String,
    #[serde(rename = "Quantity")]
    quantity: u64,
    #[serde(rename = "InvoiceDate")]
    #[serde(with = "my_date_format")]
    invoice_date: DateTime<Utc>,
    #[serde(rename = "UnitPrice")]
    unit_price: f64,
    #[serde(rename = "CustomerID")]
    customer_id: u64,
    #[serde(rename = "Country")]
    country: String,
}

mod my_date_format {
    use chrono::{DateTime, TimeZone, Utc};
    use serde::de::Error;
    use serde::{self, Deserialize, Deserializer, Serializer};

    const FORMAT: &str = "%m/%d/%Y %H:%M";

    pub fn serialize<S>(date: &DateTime<Utc>, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let s = format!("{}", date.format(FORMAT));
        serializer.serialize_str(&s)
    }

    pub fn deserialize<'de, D>(deserializer: D) -> Result<DateTime<Utc>, D::Error>
    where
        D: Deserializer<'de>,
    {
        let s = String::deserialize(deserializer)?;
        Utc.datetime_from_str(&s, FORMAT).map_err(Error::custom)
    }
}

fn map_reduce<I, K, V, Map, Red>(
    data: impl IntoParallelIterator<Item = I>,
    map: Map,
    reduce: Red,
) -> impl IntoIterator<Item = (K, V)>
where
    Map: Fn(I) -> (K, V) + Sync + Send,
    Red: Fn(V, V) -> V + Sync + Send + Copy,
    I: Sync + Send,
    K: Sync + Send + Ord + Hash,
    V: Sync + Send + Default,
{
    let data: Vec<(K, V)> = data.into_par_iter().map(map).collect();

    let mut hash: HashMap<K, Vec<V>> = HashMap::with_capacity(data.len());
    for (key, value) in data.into_iter() {
        hash.entry(key).or_insert_with(Vec::new).push(value)
    }

    let result: BTreeMap<K, V> = hash
        .into_par_iter()
        .map(|(k, v)| (k, v.into_par_iter().reduce(V::default, reduce)))
        .collect();
    result
}

fn main() -> Result<(), Box<dyn Error>> {
    let mut rdr = csv::Reader::from_reader(io::stdin());
    let map = |x: Record| (x.invoice_date.date(), x.unit_price * (x.quantity as f64));
    let reduce = |x, y| x + y;
    let data: Vec<Record> = rdr.deserialize().filter_map(|r| r.ok()).collect();

    let result = map_reduce(data, map, reduce);

    for x in result {
        println!("{} -> {}", x.0, x.1)
    }
    Ok(())
}
