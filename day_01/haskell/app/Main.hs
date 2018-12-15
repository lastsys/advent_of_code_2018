module Main where

import qualified Data.Set as Set

loadFile :: FilePath -> IO [String]
loadFile filename = fmap lines $ readFile filename

parseInt :: String -> Int
parseInt s = case s of
    (x:xs) | x == '+' -> read xs
    x                 -> read x

loadData :: FilePath -> IO [Int]
loadData filename = fmap (map parseInt) $ loadFile filename

part1 :: [Int] -> Int
part1 = sum

part2 :: [Int] -> Int
part2 values =
    fst $ head $ dropWhile condition $ scanl update (0, Set.empty) (cycle values)
    where update :: (Int, Set.Set Int) -> Int -> (Int, Set.Set Int)
          update (sum, state) value = (sum + value, Set.insert sum state)
          condition :: (Int, Set.Set Int) -> Bool
          condition (sum, state) = not $ Set.member sum state

main :: IO ()
main = do values <- loadData "../input.txt"
          putStrLn $ "Part 1: " ++ show (part1 values)
          putStrLn $ "Part 2: " ++ show (part2 values)
