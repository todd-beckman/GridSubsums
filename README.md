Subsums
=======

This project is incomplete. No functionality guaranteed. Current progress yields what is suspected (but not thoroughly tested) to be proper updating and reading.


The goal of this project is to produce an implementation of a 2D grid mxn of numbers with two features:

- Update any single cell
- Find the sum of any subrectangle within this grid

The naive solution of this problem finds the subsum in O(mn) running time. The target running time is O(log(mn)) while retaining a reasonable space complexity such as O(mnlog(mn)). A third power polynomial or worse is restricted.

