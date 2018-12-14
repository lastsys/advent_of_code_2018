// Learn more about F# at http://fsharp.org
module Program

open System

module Seq =
    let rec repeat items = 
      seq { yield! items  
            yield! repeat items }

let load_file filename =
    filename
    |> System.IO.File.ReadLines
    |> Seq.map int
    |> Seq.toArray

let part1 values = Seq.sum values

let part2 values =
    let scanner state value = (fst state + value, state ||> Set.add)
    let skipCondition state = state ||> Set.contains |> not
    Seq.repeat values
    |> Seq.scan scanner (0, Set.empty)
    |> Seq.skipWhile skipCondition
    |> Seq.head
    |> fst

[<EntryPoint>]
let main argv =
    let values = load_file "../input.txt"
    values |> part1 |> printfn "Part 1: %A"
    values |> part2 |> printfn "Part 2: %A"
    0
