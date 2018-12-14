use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;
use std::collections::{LinkedList, HashMap};

trait Sum {
    fn sum(&self) -> i32;
}

impl Sum for LinkedList<i32> {
    fn sum(&self) -> i32 {
        let mut s = 0;
        for number in self {
            s += number;
        }
        s
    }
}

fn load_file(filename: &str) -> LinkedList<i32> {
    let f = File::open(filename).expect("unable to open file");
    let reader = BufReader::new(f);
    let mut values: LinkedList<i32> = LinkedList::new();

    for line in reader.lines() {
        values.push_back(line.unwrap().parse::<i32>().unwrap());
    }

    values
}

fn part1(values: &LinkedList<i32>) {
    println!("Part 1: {}", values.sum());
}

fn part2(values: &LinkedList<i32>) {
    let mut visited = HashMap::new();
    let mut v = 0;
    loop {
        for value in values {
            v += value;
            if visited.contains_key(&v) {
                println!("Part 2: {}", v);
                return;
            } else {
                visited.insert(v, true);
            }
        }
    }
}

fn main() {
    let values = load_file("../input.txt");
    part1(&values);
    part2(&values);
}
