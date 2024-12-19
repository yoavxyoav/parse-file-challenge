# File Parsing Challenge

## General description
This challenge is about optimization, understanding your tools and, most importantly, how computers work.

## What do you need to do?
The repository contains skeleton code in the following languages: JS, Go and Zig.
Your task is to pick a language you feel comfortable with, implement the `parse` method and make it run as fast as you can.
Your implementation should read the point.txt file, sum each column and return the avg. of each column. The main file should run your implementation without any errors.

## How it goes
The repo contains a generator written in Golang that will produce 2 UTF-8 encoded files. 
The first (points.txt) is a list of floats ranging between -99.99 and 99.99 (2 decimal precision), separated by a comma ','
Lines are separated with a new line char '\n'
The second file is a verification file you can use to verify your code is correct with the same structure as the points file.
The verification file contains 2 floats and a single integer. The floats are the sum of all numbers on the same column, and the integer is the number of lines generated. The line and file ends with a new line '\n'

Each main file (the name is language dependent) contains a code to run the `parse` method, measure its time, verify its results and run in repetition till the best time is found.
Do not change any existing code; you should change only the parse method (you can change its signature and how it's being called if you want).

## How to win
Make the `parse` method as fast as you can.
I added some benchmarks for how fast it can run just so you have something to measure against. 
Whoever is fastest - wins.
If more than a single dev gets the exact timing, whoever submitted her results first - wins.
I will measure your code on my laptop precisely like I measured the benchmarks below. I will test everyone's code with the same generated points.txt file.

NOTE: Please do NOT consider the benchmark as the fastest possible implementation!!! You should try to make your code run even faster - it is not only possible but number can go down a LOT!!!

## Submitting results
Fork the repo into your account and send a PR in github. 
You may update your PR if you want but for timing purposes only the latest commit will count. 

### *Rules*
1. You can touch only the parse method; you may alter its signature, how its being called and how you handle its results.
2. You may add any additional code/methods you need.
3. You must NOT touch any of the repetition tester or verification code or alter time measurement in any way.
4. You must use only the language's internal tools.
5. You can't use any imported library.
6. You can't go outside the bands of the language - meaning no C bindings.
7. There are no other rules; you can use any trick or knowledge you have to make it work.

## Measuring time
Each implementation contains a way to measure the execution time of the parse method itself - please don't change it.
I will run each implementation with a repetition tester and take the best time from all the runs.
Each implementation contains a readme file with building instructions, and these are the build instructions I will be using.

## Some hints and tips
1. First, generate a small file, which can be easier to work with; 100,000 lines is a good place to start.
2. Sometimes, you will need to debug the code, so generating a file with 10 lines (or even 2) might help.
3. There are no rules against using chatgpt/LLM or searching the web, BUT I highly recommend you try to figure it out yourself.
   Learning how to measure performance and optimize parts of code is a skill you should pick up and not trust gpitpot to do for you.
4. measure, measure and finally measure; that's the key for advancing.

## Bonus
I will give a bonus point if you can a write a function that finds the fastest theoretical (!!!) run time possible.
Besides a bonus, it can help you understand what's possible and what's not.

## Generating numbers
Install Go with brew or from the Golang website - use 1.23.2 (use this also in case you want to solve the Golang challenge)

```zsh
cd generator
go build -ldflags "-s -w" -o=points-generator main.go
./points-generator
```
Standing inside an implementation directory you can try:
`../points-generator/points-generator`

## Results
here are some benchmark for 100,000,000 points (~= 1.2G) on both Apple M4 (will be used score results) and M1.

Js benchmark (node) - 1.2s (M4), 2.1s (M1)  
Golang benchmark - 1s (M4), 1.6s (M1)  
Zig benchmark - 852ms (M4), 1.1s (M1)  

Good luck and have fun!!!

