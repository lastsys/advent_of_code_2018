use std::collections::HashMap;
use std::iter::Iterator;
use std::fs::File;
use std::io::BufReader;
use std::io::prelude::*;

fn load_file(filename: &str) -> Vec<String> {
    let f = File::open(filename).expect("unable to open file");
    let reader = BufReader::new(f);
    let mut serials: Vec<String> = Vec::new();

    for line in reader.lines() {
        serials.push(line.unwrap());
    }

    serials
}

fn count_pairs_and_triplets(s: &String) -> (u32, u32) {
    let mut count: HashMap<char, u32> = HashMap::new();
    for c in s.chars() {
        let char_count = count.entry(c).or_insert(0);
        *char_count += 1
    }

    count.values()
        .fold((0u32, 0u32),
              |state, val|
                  match val {
                      2 => (state.0 + 1, state.1),
                      3 => (state.0, state.1 + 1),
                      _ => state
                  })
}

fn part1(serials: &Vec<String>) {
    let mut pairs = 0;
    let mut triplets = 0;
    for serial in serials {
        let (p, t) = count_pairs_and_triplets(serial);
        if p >= 1 {
            pairs += 1;
        }
        if t >= 1 {
            triplets += 1;
        }
    }
    let checksum = pairs * triplets;
    println!("Part1: {} * {} = {}", pairs, triplets, checksum);
}

fn number_of_non_equal_characters(s1: &String, s2: &String) -> (u32, String) {
    let mut rest = String::new();
    let mut count = 0u32;
    for (c1, c2) in s1.chars().zip(s2.chars()) {
        if c1 != c2 {
            count += 1;
        } else {
            rest.push(c1);
        }
    }
    (count, rest)
}

fn part2(serials: &Vec<String>) {
    for (i, serial1) in serials.iter().enumerate() {
        for (j, serial2) in serials.iter().enumerate() {
            if j > i && serial1 != serial2 {
                let (count, rest) = number_of_non_equal_characters(serial1, serial2);
                if count == 1 {
                    println!("{}, {}", count, rest);
                }
            }
        }
    }
}

fn main() {
    let serials = load_file("../input.txt");
    part1(&serials);
    part2(&serials);
}
