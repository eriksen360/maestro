package lib

/*
1) Rewrite parser lib so that it can handle files at random
2) Write parser lib more elegantly. A single file should be able to define both hosts and tasks as
top-level keys
3) Write validation for yaml scripts

4) Write a CLI interface that consumes the .yaml files and stores them in such a fashion that
the executor can easily consume them and look ids, tasks and instructions up (as Tasks)
Topology:

		 User
       /     \
      /       \
   Parser----Executor



*/
