# Duplicate File Handler

## About 
Duplicate File Handler is a useful tool that can free some space on your drive. 
The handler
- checks and compares files in a folder, 
- displays the result, and 
- removes duplicates.

## Learning outcomes
- Learn how to work with files and folders. 
- Get familiar with hashing, learn how to apply it to your tasks.

## Stages
1. Scan files and folders with the help of the filepath package.
2. Use file sizes to find duplicate files.
3. Learn about hash functions and implement them in your code.
4. Let's delete all duplicates.

## Objectives
1. Accept a command-line argument that is a root directory with files and folders. 
2. Output the message Directory is not specified if there is no command-line argument.
3. Read user input that specifies the file format. Empty input should match any file format.
4. Print a menu with two sorting options: Descending and Ascending.
5. Iterate over folders and print the information about files of the same size: their size, path, and names.
6. Ask for duplicates check. If the input is yes, get the hash of files of the same size, group the files of the same hash, and assign numbers to these files.
7. Ask a user whether they want to delete files. If yes, read what files a user wants to delete and then delete them.
8. Read a sequence of files that a user wants to delete and then delete them. 
9. Print the total freed-up space in bytes.

## Example
Suppose you have the following set of files and folders:
```go
+---[root_folder]
    +---gordon_ramsay_chicken_breast.avi /4590560 bytes
    +---poker_face.mp3 /5550640 bytes
    +---poker_face_copy.mp3 /5550640 bytes
    +---[audio]
    |   |
    |   +---voice.mp3 /2319746 bytes
    |   +---sia_snowman.mp3 /4590560 bytes
    |   +---nea_some_say.mp3 /3232056 bytes
    |   +---[classic]
    |   |   |
    |   |   +---unknown.mp3 /3422208 bytes
    |   |   +---vivaldi_four_seasons_winter.mp3 /9158144 bytes
    |   |   +---chopin_waltz7_op64_no2.mp3 /9765504 bytes
    |   +---[rock]
    |       |
    |       +---smells_like_teen_spirit.mp3 /4590560 bytes
    |       +---numb.mp3 /5786312 bytes
    +---[masterpiece]
        |
        +---rick_astley_never_gonna_give_you_up.mp3 /3422208 bytes
```
Program output:
```go
> go run main.go root_folder

Enter file format:
>

Size sorting options:
1. Descending
2. Ascending

Enter a sorting option:
> 1

5550640 bytes
root_folder/poker_face.mp3
root_folder/poker_face_copy.mp3

4590560 bytes
root_folder/gordon_ramsay_chicken_breast.avi
root_folder/audio/sia_snowman.mp3
root_folder/audio/rock/smells_like_teen_spirit.mp3

3422208 bytes
root_folder/audio/classic/unknown.mp3
root_folder/masterpiece/rick_astley_never_gonna_give_you_up.mp3

Check for duplicates?
> yes

5550640 bytes
Hash: 909ba4ad2bda46b10aac3c5b7f01abd5
1. root_folder/poker_face.mp3
2. root_folder/poker_face_copy.mp3

3422208 bytes
Hash: a7f5f35426b927411fc9231b56382173
3. root_folder/audio/classic/unknown.mp3
4. root_folder/masterpiece/rick_astley_never_gonna_give_you_up.mp3

Delete files?
> yes

Enter file numbers to delete:
> 1 2 5

Wrong format

Enter file numbers to delete:
> 1 2 4

Total freed up space: 14523488 bytes
```