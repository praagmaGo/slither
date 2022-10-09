# slither.io robot

This is the server part of a program who should play alone at the game slither.io

**What it is**
A first functional version was previously written with webworkers, directly on modifying the game web page
Finally I choose to write a server to be able to manage several users from a central place. They are intended to collaborate.
The main technical challenge is done: sending all the game infos to the server (using javascript directly on the original webpage). Then all the data is distributed on local structures. All this take always less then 1 ms.

**Why golang**
I choose to write the server with go. This program is doing all analysis and decisions.
Go is compiled, and has very good possibilities for managing several processes, thus allowing to more strategies on less time.
(See https://gobyexample.com/worker-pools site)
This code is finally intended to fight against real good players or other robots.
This code is written to be (very) fast. You can have a look at the function GetIdx for an example.
