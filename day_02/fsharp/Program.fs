open System

/// <summary>Load input data.</summary>
/// <param name="filename">Name of file to load.</param>
/// <returns>Sequence of strings, one per row in file.</returns>
let loadFile filename =
    System.IO.File.ReadLines(filename)
    |> Seq.map (fun s -> s.Trim())

/// <summary>Count pairs and triplets in a single string.</summary>
/// <param name="s">String to process.</param>
/// <returns>Tuple with number of pairs and triplets.</returns>
let countPairsAndTriplets s =
    let countChars =
        fun count char -> match (Map.tryFind char count) with
                          | Some(x) -> Map.add char (x + 1) count  // increase character count
                          | None    -> Map.add char 1 count        // if character has not been seen
    let countTwoAndThree =
        fun c k v -> match v with
                     | v when v = 2 -> ((fst c) + 1, snd c)  // increase pair count
                     | v when v = 3 -> (fst c, (snd c) + 1)  // increase triplet count
                     | _            -> c                     // do nothing
    s |> Seq.fold countChars Map.empty |> Map.fold countTwoAndThree (0, 0)

/// <summary>Part 1 implementation.</summary>
/// <param name="serials">Sequence of serials.</param>
let part1 serials =
    let countThem count serial =
        let (p, t) = countPairsAndTriplets serial
        ((if p >= 1 then (fst count) + 1 else (fst count)),
         (if t >= 1 then (snd count) + 1 else (snd count)))
                
    let (pairs, triplets) = serials |> Seq.fold countThem (0, 0)
    let checksum = pairs * triplets
    printfn "Part1: %d * %d = %d" pairs triplets checksum
    
/// <summary>Count number of non-equal characters in two strings.<summary>
/// <param name="s1">First string.</param>
/// <param name="s2">Second string.</param>
/// <returns>Tuple with number of non-equal characters and a list of equal characters.</returns>
let numberOfNonEqualCharacters (s1: string) (s2: string) =
    let f (state: int*List<char>) (pair: char*char) =
        if ((fst pair) <> (snd pair))
        then ((fst state) + 1, snd state)
        else (fst state, (fst pair) :: (snd state))
    Seq.zip s1 s2 |> Seq.fold f (0, List.empty)
    
/// <summary>Part 2 implementation.</summary>
/// <param name="serials">Sequence of serials.</param>
let part2 (serials: seq<string>) =       
    let r = seq { for (i, s1) in Seq.mapi (fun i s -> (i, s)) serials do
                      for s2 in Seq.skip i serials do
                          if s1 <> s2 then
                              let (count, rest) = numberOfNonEqualCharacters s1 s2
                              if count = 1 then
                                  yield (count, Array.ofList rest |> Array.rev |> String.Concat) }
    printfn "Part2: %A" r
    
[<EntryPoint>]
let main argv =
    let serials = loadFile "..\input.txt"
    part1 serials
    part2 serials
    0 // return an integer exit code
